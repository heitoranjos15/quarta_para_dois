package nflverse

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"time"
)

const gameIndexTTL = 24 * time.Hour

type GameRef struct {
	GameID, ESPNID, GameType, HomeTeam, AwayTeam, Roof, Surface, StadiumID, Stadium string
	Season, Week                                                                    int
}

type GameIndex struct {
	byGameID map[string]GameRef
	byESPNID map[string]GameRef
	// gameType of any game in a (season, week); all games of a week share a ESPN week mapping
	weekType map[[2]int]string
}

// ParseGamesCSV reads nflverse's games.csv, locating columns by header name.
func ParseGamesCSV(r io.Reader) (*GameIndex, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read games.csv header: %w", err)
	}
	col := make(map[string]int, len(header))
	for i, h := range header {
		col[h] = i
	}
	for _, need := range []string{"game_id", "season", "game_type", "week", "home_team", "away_team"} {
		if _, ok := col[need]; !ok {
			return nil, fmt.Errorf("games.csv: missing column %q", need)
		}
	}

	ix := &GameIndex{
		byGameID: map[string]GameRef{},
		byESPNID: map[string]GameRef{},
		weekType: map[[2]int]string{},
	}
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read games.csv: %w", err)
		}
		get := func(name string) string {
			if i, ok := col[name]; ok && i < len(rec) {
				return rec[i]
			}
			return ""
		}
		season, errS := strconv.Atoi(get("season"))
		week, errW := strconv.Atoi(get("week"))
		if errS != nil || errW != nil || get("game_id") == "" {
			continue
		}
		ref := GameRef{
			GameID: get("game_id"), ESPNID: get("espn"), GameType: get("game_type"),
			HomeTeam: get("home_team"), AwayTeam: get("away_team"),
			Roof: get("roof"), Surface: get("surface"),
			StadiumID: get("stadium_id"), Stadium: get("stadium"),
			Season: season, Week: week,
		}
		ix.byGameID[ref.GameID] = ref
		if ref.ESPNID != "" {
			ix.byESPNID[ref.ESPNID] = ref
		}
		ix.weekType[[2]int{season, week}] = ref.GameType
	}
	return ix, nil
}

func (ix *GameIndex) ByGameID(id string) (GameRef, bool) {
	ref, ok := ix.byGameID[id]
	return ref, ok
}

func (ix *GameIndex) ByESPNID(id string) (GameRef, bool) {
	ref, ok := ix.byESPNID[id]
	return ref, ok
}

func (ix *GameIndex) Seasons() []int {
	seen := map[int]bool{}
	for k := range ix.weekType {
		seen[k[0]] = true
	}
	out := make([]int, 0, len(seen))
	for s := range seen {
		out = append(out, s)
	}
	sort.Ints(out)
	return out
}

// Weeks lists a season's weeks in ascending order, playoff weeks included.
func (ix *GameIndex) Weeks(season int) []int {
	var out []int
	for k := range ix.weekType {
		if k[0] == season {
			out = append(out, k[1])
		}
	}
	sort.Ints(out)
	return out
}

// ESPNWeek maps an nflverse (season, week) to ESPN's seasontype and week number.
func (ix *GameIndex) ESPNWeek(season, week int) (seasonType, espnWeek int, ok bool) {
	switch ix.weekType[[2]int{season, week}] {
	case "REG":
		return 2, week, true
	case "WC":
		return 3, 1, true
	case "DIV":
		return 3, 2, true
	case "CON":
		return 3, 3, true
	case "SB":
		return 3, 5, true // ESPN skips week 4 (Pro Bowl)
	}
	return 0, 0, false
}

// GameIndex returns the games.csv index, held in memory for 24 hours. If a refresh
// fails, the previous index keeps being served.
func (c *Client) GameIndex(ctx context.Context) (*GameIndex, error) {
	c.mu.Lock()
	ix, at := c.gameIndex, c.gameIndexAt
	c.mu.Unlock()
	if ix != nil && c.opts.Now().Sub(at) < gameIndexTTL {
		return ix, nil
	}

	res, err, _ := c.sfGroup.Do("games-index", func() (any, error) {
		data, err := c.download(context.WithoutCancel(ctx), c.opts.GamesCSVURL)
		if err != nil {
			return nil, err
		}
		return ParseGamesCSV(bytes.NewReader(data))
	})
	if err != nil {
		if ix != nil {
			slog.Warn("refresh games.csv; serving previous index", "err", err)
			return ix, nil
		}
		return nil, err
	}

	fresh := res.(*GameIndex)
	c.mu.Lock()
	c.gameIndex, c.gameIndexAt = fresh, c.opts.Now()
	c.mu.Unlock()
	return fresh, nil
}

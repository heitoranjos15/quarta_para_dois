// Package games joins ESPN's schedule and live scores to nflverse's game index.
package games

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"quarta-para-dois/backend/internal/espn"
	"quarta-para-dois/backend/internal/nflverse"
)

var ErrNotFound = errors.New("game not found")

type Game struct {
	GameID    string `json:"game_id"`
	ESPNID    string `json:"espn_id"`
	Season    int    `json:"season"`
	Week      int    `json:"week"`
	GameType  string `json:"game_type"`
	Kickoff   string `json:"kickoff"` // RFC3339 UTC, "" if unknown
	Status    string `json:"status"`  // scheduled | in_progress | final
	HomeTeam  string `json:"home_team"`
	AwayTeam  string `json:"away_team"`
	HomeScore *int   `json:"home_score"`
	AwayScore *int   `json:"away_score"`
	Stadium   string `json:"stadium"`
	Location  string `json:"location"`
	Indoor    bool   `json:"indoor"`
	Roof      string `json:"roof,omitempty"`
	Surface   string `json:"surface,omitempty"`
}

type IndexSource interface {
	GameIndex(ctx context.Context) (*nflverse.GameIndex, error)
}

type ESPNSource interface {
	Scoreboard(ctx context.Context, season, seasonType, week int) ([]espn.Event, error)
	TeamSchedule(ctx context.Context, nflverseAbbr string, season int) ([]espn.Event, error)
}

type Service struct {
	ix IndexSource
	es ESPNSource
}

func NewService(ix IndexSource, es ESPNSource) *Service {
	return &Service{ix: ix, es: es}
}

func (s *Service) Seasons(ctx context.Context) ([]int, error) {
	ix, err := s.ix.GameIndex(ctx)
	if err != nil {
		return nil, err
	}
	return ix.Seasons(), nil
}

func (s *Service) Weeks(ctx context.Context, season int) ([]int, error) {
	ix, err := s.ix.GameIndex(ctx)
	if err != nil {
		return nil, err
	}
	return ix.Weeks(season), nil
}

func (s *Service) Week(ctx context.Context, season, week int) ([]Game, error) {
	ix, err := s.ix.GameIndex(ctx)
	if err != nil {
		return nil, err
	}
	events, err := s.weekEvents(ctx, ix, season, week)
	if err != nil {
		return nil, err
	}
	return toGames(ix, events, season, week), nil
}

// Game finds one game with a single scoreboard request for its week.
func (s *Service) Game(ctx context.Context, gameID string) (Game, error) {
	ix, err := s.ix.GameIndex(ctx)
	if err != nil {
		return Game{}, err
	}

	if ref, ok := ix.ByGameID(gameID); ok {
		if ref.ESPNID == "" {
			return Game{}, ErrNotFound
		}
		events, err := s.weekEvents(ctx, ix, ref.Season, ref.Week)
		if err != nil {
			return Game{}, err
		}
		for _, e := range events {
			if e.ID == ref.ESPNID {
				return toGame(ix, e, ref.Season, ref.Week), nil
			}
		}
		return Game{}, ErrNotFound
	}

	// Not in games.csv yet: it may still be a game ESPN lists under a fallback ID.
	season, week, ok := nflverse.ParseGameID(gameID)
	if !ok {
		return Game{}, ErrNotFound
	}
	events, err := s.weekEvents(ctx, ix, season, week)
	if err != nil {
		return Game{}, err
	}
	for _, g := range toGames(ix, events, season, week) {
		if g.GameID == gameID {
			return g, nil
		}
	}
	return Game{}, ErrNotFound
}

func (s *Service) TeamGames(ctx context.Context, abbr string, season int) ([]Game, error) {
	ix, err := s.ix.GameIndex(ctx)
	if err != nil {
		return nil, err
	}
	events, err := s.es.TeamSchedule(ctx, abbr, season)
	if err != nil {
		return nil, err
	}
	out := make([]Game, 0, len(events))
	for _, e := range events {
		out = append(out, toGame(ix, e, season, e.Week))
	}
	return out, nil
}

// weekEvents loads one week's ESPN scoreboard. A regular-season week that games.csv
// doesn't list yet is still asked for, so a newly published week isn't empty.
func (s *Service) weekEvents(ctx context.Context, ix *nflverse.GameIndex, season, week int) ([]espn.Event, error) {
	seasonType, espnWeek, ok := ix.ESPNWeek(season, week)
	if !ok {
		if week < 1 || week > 18 {
			return nil, ErrNotFound
		}
		seasonType, espnWeek = 2, week
	}
	return s.es.Scoreboard(ctx, season, seasonType, espnWeek)
}

func toGames(ix *nflverse.GameIndex, events []espn.Event, season, week int) []Game {
	out := make([]Game, 0, len(events))
	for _, e := range events {
		out = append(out, toGame(ix, e, season, week))
	}
	return out
}

// toGame builds a Game from an ESPN event. When games.csv knows the event, the IDs,
// teams, week and type come from it; otherwise ESPN's data is used as a fallback.
func toGame(ix *nflverse.GameIndex, e espn.Event, season, week int) Game {
	g := Game{
		ESPNID:    e.ID,
		Season:    season,
		Week:      week,
		Status:    string(e.Status),
		HomeScore: e.HomeScore,
		AwayScore: e.AwayScore,
		Stadium:   e.Venue,
		Location:  location(e.City, e.State),
		Indoor:    e.Indoor,
	}
	if !e.Kickoff.IsZero() {
		g.Kickoff = e.Kickoff.UTC().Format(time.RFC3339)
	}

	if ref, ok := ix.ByESPNID(e.ID); ok {
		g.GameID, g.Season, g.Week, g.GameType = ref.GameID, ref.Season, ref.Week, ref.GameType
		g.HomeTeam, g.AwayTeam = ref.HomeTeam, ref.AwayTeam
		g.Roof, g.Surface = ref.Roof, ref.Surface
		if g.Stadium == "" {
			g.Stadium = ref.Stadium
		}
		return g
	}

	if e.Week > 0 {
		g.Week = e.Week
	}
	g.HomeTeam, g.AwayTeam = normalizeAbbr(e.HomeAbbr), normalizeAbbr(e.AwayAbbr)
	g.GameID = fmt.Sprintf("%d_%02d_%s_%s", g.Season, g.Week, g.AwayTeam, g.HomeTeam)
	return g
}

// normalizeAbbr converts ESPN team abbreviations to nflverse's.
func normalizeAbbr(abbr string) string {
	switch a := strings.ToUpper(abbr); a {
	case "WSH":
		return "WAS"
	case "LAR":
		return "LA"
	default:
		return a
	}
}

func location(city, state string) string {
	switch {
	case city != "" && state != "":
		return city + ", " + state
	case city != "":
		return city
	}
	return state
}

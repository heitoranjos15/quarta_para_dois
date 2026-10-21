package nflverse

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var ErrGameNotFound = errors.New("game not found")

const gamePlaysTTL = 30 * 24 * time.Hour

type Options struct {
	PBPURLFormat  string        // default "https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_%d.parquet"
	GamesCSVURL   string        // default "https://github.com/nflverse/nfldata/raw/master/data/games.csv"
	MemoryGames   int           // default 64
	SeasonRefetch time.Duration // default 1h: min gap between season downloads
	Now           func() time.Time
}

func (o Options) withDefaults() Options {
	if o.PBPURLFormat == "" {
		o.PBPURLFormat = "https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_%d.parquet"
	}
	if o.GamesCSVURL == "" {
		o.GamesCSVURL = "https://github.com/nflverse/nfldata/raw/master/data/games.csv"
	}
	if o.MemoryGames <= 0 {
		o.MemoryGames = 64
	}
	if o.SeasonRefetch <= 0 {
		o.SeasonRefetch = time.Hour
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	return o
}

func gamePlaysKey(gameID string) string { return "pbp:game:" + gameID }

// GamePlays returns one game's plays. Lookup order: memory, Redis, then a download
// of the whole season, which is split per game, written to Redis and dropped.
func (c *Client) GamePlays(ctx context.Context, gameID string) ([]Play, error) {
	season, _, ok := ParseGameID(gameID)
	if !ok {
		return nil, fmt.Errorf("invalid game id %q", gameID)
	}

	if plays, ok := c.games.Get(gameID); ok {
		return plays, nil
	}

	var plays []Play
	hit, err := c.cache.Get(ctx, gamePlaysKey(gameID), &plays)
	if err != nil {
		slog.Warn("cache get", "key", gamePlaysKey(gameID), "err", err)
	}
	if hit {
		c.games.Add(gameID, plays)
		return plays, nil
	}

	if c.recentlyMissing(season, gameID) {
		return nil, ErrGameNotFound
	}

	result, err, _ := c.sfGroup.Do(fmt.Sprintf("pbp-season:%d", season), func() (any, error) {
		// Shared by every waiter, so one caller giving up must not cancel it.
		return c.downloadSeason(context.WithoutCancel(ctx), season)
	})
	if err != nil {
		return nil, err
	}

	plays = result.(map[string][]Play)[gameID]
	if len(plays) == 0 {
		return nil, ErrGameNotFound
	}
	c.games.Add(gameID, plays)
	return plays, nil
}

// seasonState remembers the last download of a season.
type seasonState struct {
	fetchedAt time.Time
	inFile    map[string]bool // games present in that file
	stored    map[string]int  // games written to Redis, by play count
}

// recentlyMissing is true for a game that wasn't in a download made less than SeasonRefetch
// ago, so requests for a game nflverse hasn't published can't each trigger a download.
// A game that was in the file is never throttled: if it isn't cached, Redis lost it.
func (c *Client) recentlyMissing(season int, gameID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	st, ok := c.seasons[season]
	return ok && !st.inFile[gameID] && c.opts.Now().Sub(st.fetchedAt) < c.opts.SeasonRefetch
}

// downloadSeason fetches a season's play-by-play and writes to Redis the games that are
// new or changed since the last write.
func (c *Client) downloadSeason(ctx context.Context, season int) (map[string][]Play, error) {
	data, err := c.download(ctx, fmt.Sprintf(c.opts.PBPURLFormat, season))
	if err != nil {
		return nil, err
	}
	all, err := ParsePlays(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	byGame := make(map[string][]Play)
	for _, p := range all {
		byGame[p.GameID] = append(byGame[p.GameID], p)
	}

	c.mu.Lock()
	st, ok := c.seasons[season]
	if !ok {
		st = &seasonState{stored: map[string]int{}}
		c.seasons[season] = st
	}
	st.fetchedAt = c.opts.Now()
	st.inFile = make(map[string]bool, len(byGame))
	toWrite := map[string][]Play{}
	for id, plays := range byGame {
		st.inFile[id] = true
		if st.stored[id] != len(plays) {
			toWrite[id] = plays
		}
	}
	c.mu.Unlock()

	for id, plays := range toWrite {
		if err := c.cache.Set(ctx, gamePlaysKey(id), plays, gamePlaysTTL); err != nil {
			slog.Warn("cache set", "key", gamePlaysKey(id), "err", err)
			continue
		}
		c.mu.Lock()
		st.stored[id] = len(plays)
		c.mu.Unlock()
	}
	return byGame, nil
}

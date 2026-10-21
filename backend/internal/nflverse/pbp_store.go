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

// downloadSeason fetches a season's play-by-play and stores every game in Redis.
// A season fetched less than SeasonRefetch ago is not fetched again, so requests
// for a game that isn't in the file can't trigger a download each time.
func (c *Client) downloadSeason(ctx context.Context, season int) (map[string][]Play, error) {
	c.mu.Lock()
	last, fetched := c.seasonFetchedAt[season]
	c.mu.Unlock()
	if fetched && c.opts.Now().Sub(last) < c.opts.SeasonRefetch {
		return nil, ErrGameNotFound
	}

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
	for id, plays := range byGame {
		if err := c.cache.Set(ctx, gamePlaysKey(id), plays, gamePlaysTTL); err != nil {
			slog.Warn("cache set", "key", gamePlaysKey(id), "err", err)
		}
	}

	c.mu.Lock()
	c.seasonFetchedAt[season] = c.opts.Now()
	c.mu.Unlock()
	return byGame, nil
}

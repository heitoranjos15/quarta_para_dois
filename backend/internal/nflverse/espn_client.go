package nflverse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"quarta-para-dois/backend/internal/cache"
)

const (
	espnBaseURL  = "https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard"
	espnCacheTTL = 24 * time.Hour
	espnCacheKey = "espn:schedule:%d:%d" // season:week
)

type ESPNGamesResponse struct {
	Leagues []ESPNLeague `json:"leagues"`
	Events  []ESPNEvent  `json:"events"`
	Season  ESPNSeason   `json:"season"`
	Week    struct {
		Number int `json:"number"`
	} `json:"week"`
}

type ESPNSeason struct {
	Year int         `json:"year"`
	Type interface{} `json:"type"`
}

type ESPNLeague struct {
	ID     string `json:"id"`
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Season struct {
		Year int      `json:"year"`
		Type ESPNType `json:"type"`
	} `json:"season"`
}

type ESPNType struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Name string `json:"name"`
}

type ESPNTeamScheduleResponse struct {
	Events []ESPNEvent `json:"events"`
	Season struct {
		Year int `json:"year"`
	} `json:"season"`
	Team struct {
		Abbreviation string `json:"abbreviation"`
	} `json:"team"`
}

type ESPNEvent struct {
	ID           string            `json:"id"`
	Date         string            `json:"date"`
	Name         string            `json:"name"`
	ShortName    string            `json:"shortName"`
	Competitions []ESPNCompetition `json:"competitions"`
	Status       ESPNStatus        `json:"status"`
}

type ESPNCompetition struct {
	ID          string           `json:"id"`
	Date        string           `json:"date"`
	Status      ESPNStatus       `json:"status"`
	Venue       ESPNVenue        `json:"venue"`
	Competitors []ESPNCompetitor `json:"competitors"`
}

type ESPNStatus struct {
	Type ESPNStatusType `json:"type"`
}

type ESPNStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	ShortDetail string `json:"shortDetail"`
}

type ESPNVenue struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Address  struct {
		City  string `json:"city"`
		State string `json:"state"`
	} `json:"address"`
}

type ESPNCompetitor struct {
	ID       string      `json:"id"`
	Team     ESPNTeam    `json:"team"`
	Score    interface{} `json:"score"`
	HomeAway string      `json:"homeAway"`
	Winner   bool        `json:"winner"`
}

func (c *ESPNCompetitor) GetScoreValue() float64 {
	switch v := c.Score.(type) {
	case string:
		var val float64
		_, _ = fmt.Sscanf(v, "%f", &val)
		return val
	case map[string]interface{}:
		if val, ok := v["value"].(float64); ok {
			return val
		}
	case ESPNScore:
		return v.Value
	}
	return 0
}

type ESPNScore struct {
	Value        float64 `json:"value"`
	DisplayValue string  `json:"displayValue"`
}

type ESPNTeam struct {
	ID               string `json:"id"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Logo             string `json:"logo"`
	Color            string `json:"color"`
}

type ESPNClient struct {
	httpClient *http.Client
	cache      cache.Cache
}

func NewESPNClient(c cache.Cache) *ESPNClient {
	return &ESPNClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      c,
	}
}

func (c *ESPNClient) GetSchedule(ctx context.Context, season, week int) ([]Game, error) {
	cacheKey := fmt.Sprintf(espnCacheKey, season, week)

	// Try cache first
	var cached []Game
	if ok, err := c.cache.Get(ctx, cacheKey, &cached); err != nil {
		slog.Warn("cache get", "key", cacheKey, "err", err)
	} else if ok {
		return cached, nil
	}

	// Build ESPN URL
	url := fmt.Sprintf("%s?seasontype=2&week=%d&dates=%d", espnBaseURL, week, season)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ESPN API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var espnResp ESPNGamesResponse
	if err := json.Unmarshal(body, &espnResp); err != nil {
		return nil, err
	}

	games := c.convertEvents(espnResp.Events, season, week)

	// Cache the result
	if err := c.cache.Set(ctx, cacheKey, games, espnCacheTTL); err != nil {
		slog.Warn("cache set", "key", cacheKey, "err", err)
	}

	return games, nil
}

func (c *ESPNClient) GetGame(ctx context.Context, season, week int, gameID string) (*Game, error) {
	if week == 0 {
		// Search all weeks (1-18)
		for w := 1; w <= 18; w++ {
			game, err := c.GetGame(ctx, season, w, gameID)
			if err != nil {
				return nil, err
			}
			if game != nil {
				return game, nil
			}
		}
		return nil, nil
	}

	games, err := c.GetSchedule(ctx, season, week)
	if err != nil {
		return nil, err
	}
	for _, g := range games {
		if g.GameID == gameID {
			return &g, nil
		}
	}
	return nil, nil
}

func (c *ESPNClient) convertEvents(events []ESPNEvent, season, week int) []Game {
	var games []Game
	for _, e := range events {
		if len(e.Competitions) == 0 {
			continue
		}
		comp := e.Competitions[0]

		var homeTeam, awayTeam string
		var homeScore, awayScore int

		for _, comp := range comp.Competitors {
			switch comp.HomeAway {
			case "home":
				homeTeam = comp.Team.Abbreviation
				homeScore = int(comp.GetScoreValue())
			case "away":
				awayTeam = comp.Team.Abbreviation
				awayScore = int(comp.GetScoreValue())
			}
		}

		// Generate game_id in our format: season_week_away_home
		gameID := fmt.Sprintf("%d_%02d_%s_%s", season, week, awayTeam, homeTeam)

		game := Game{
			GameID:    gameID,
			Season:    season,
			Week:      week,
			GameDate:  comp.Date[:10],
			HomeTeam:  homeTeam,
			AwayTeam:  awayTeam,
			HomeScore: homeScore,
			AwayScore: awayScore,
			Stadium:   comp.Venue.FullName,
			Location:  fmt.Sprintf("%s, %s", comp.Venue.Address.City, comp.Venue.Address.State),
		}
		games = append(games, game)
	}
	return games
}

func (c *ESPNClient) GetTeamSchedule(ctx context.Context, teamAbbr string, season int) ([]Game, error) {
	cacheKey := fmt.Sprintf("espn:team:%s:%d", strings.ToLower(teamAbbr), season)

	// Try cache first
	var cached []Game
	if ok, err := c.cache.Get(ctx, cacheKey, &cached); err != nil {
		slog.Warn("cache get", "key", cacheKey, "err", err)
	} else if ok {
		return cached, nil
	}

	// ESPN team abbreviation mapping
	espnAbbr := map[string]string{
		"ARI": "ari", "ATL": "atl", "BAL": "bal", "BUF": "buf",
		"CAR": "car", "CHI": "chi", "CIN": "cin", "CLE": "cle",
		"DAL": "dal", "DEN": "den", "DET": "det", "GB": "gb",
		"HOU": "hou", "IND": "ind", "JAX": "jax", "KC": "kc",
		"LV": "lv", "LAC": "lac", "LAR": "lar", "MIA": "mia",
		"MIN": "min", "NE": "ne", "NO": "no", "NYG": "nyg",
		"NYJ": "nyj", "PHI": "phi", "PIT": "pit", "SF": "sf",
		"SEA": "sea", "TB": "tb", "TEN": "ten", "WSH": "wsh",
	}

	espnTeam, ok := espnAbbr[strings.ToUpper(teamAbbr)]
	if !ok {
		return []Game{}, nil
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/%s/schedule?season=%d", espnTeam, season)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ESPN API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var teamResp ESPNTeamScheduleResponse
	if err := json.Unmarshal(body, &teamResp); err != nil {
		return nil, err
	}

	allEvents := teamResp.Events
	games := c.convertEvents(allEvents, season, 0)

	// Filter to only include games for the specified season
	var filtered []Game
	for _, g := range games {
		if g.Season == season {
			filtered = append(filtered, g)
		}
	}

	// Cache the result
	if err := c.cache.Set(ctx, cacheKey, filtered, espnCacheTTL); err != nil {
		slog.Warn("cache set", "key", cacheKey, "err", err)
	}

	return filtered, nil
}

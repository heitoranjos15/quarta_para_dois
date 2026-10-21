// Package espn reads ESPN's public NFL scoreboard and team schedule APIs.
package espn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"quarta-para-dois/backend/internal/cache"
)

const (
	defaultBaseURL = "https://site.api.espn.com/apis/site/v2/sports/football/nfl"

	finalTTL    = 30 * 24 * time.Hour
	upcomingTTL = time.Hour
	recheckTTL  = time.Minute // a game past kickoff that ESPN still lists as scheduled

	// Bump when the cached value's shape changes, so values from older builds are ignored.
	cacheVersion = "v2"
)

type Status string

const (
	StatusScheduled  Status = "scheduled"
	StatusInProgress Status = "in_progress"
	StatusFinal      Status = "final"
)

type Event struct {
	ID                 string
	Week               int // ESPN's week number, 0 if absent
	Kickoff            time.Time
	Status             Status
	HomeAbbr, AwayAbbr string // ESPN abbreviations as returned
	HomeScore          *int   // nil unless in progress or final
	AwayScore          *int
	Venue, City, State string
	Indoor             bool
}

type Options struct {
	BaseURL string
	HTTP    *http.Client
	Now     func() time.Time // default time.Now
}

type Client struct {
	cache cache.Cache
	http  *http.Client
	base  string
	now   func() time.Time
}

func NewClient(c cache.Cache, opts Options) *Client {
	if opts.BaseURL == "" {
		opts.BaseURL = defaultBaseURL
	}
	if opts.HTTP == nil {
		opts.HTTP = &http.Client{Timeout: 10 * time.Second}
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	return &Client{cache: c, http: opts.HTTP, base: strings.TrimRight(opts.BaseURL, "/"), now: opts.Now}
}

func (c *Client) Scoreboard(ctx context.Context, season, seasonType, week int) ([]Event, error) {
	key := fmt.Sprintf("espn:%s:scoreboard:%d:%d:%d", cacheVersion, season, seasonType, week)
	url := fmt.Sprintf("%s/scoreboard?seasontype=%d&week=%d&dates=%d", c.base, seasonType, week, season)
	return c.events(ctx, key, url)
}

// TeamSchedule takes an nflverse team abbreviation.
func (c *Client) TeamSchedule(ctx context.Context, nflverseAbbr string, season int) ([]Event, error) {
	slug := teamSlug(nflverseAbbr)
	key := fmt.Sprintf("espn:%s:team:%s:%d", cacheVersion, slug, season)
	url := fmt.Sprintf("%s/teams/%s/schedule?season=%d", c.base, slug, season)
	return c.events(ctx, key, url)
}

func (c *Client) events(ctx context.Context, key, url string) ([]Event, error) {
	var cached []Event
	if ok, err := c.cache.Get(ctx, key, &cached); err != nil {
		slog.Warn("cache get", "key", key, "err", err)
	} else if ok {
		return cached, nil
	}

	events, err := c.fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	if ttl, cacheable := scoreboardTTL(events, c.now()); cacheable {
		if err := c.cache.Set(ctx, key, events, ttl); err != nil {
			slog.Warn("cache set", "key", key, "err", err)
		}
	}
	return events, nil
}

// scoreboardTTL: finished games are kept 30 days. Upcoming ones are kept 1 hour, or only
// until the next kickoff if that is sooner. Nothing is cached while any game is in progress,
// and an empty answer isn't cached either (ESPN may not have published that week yet).
func scoreboardTTL(events []Event, now time.Time) (ttl time.Duration, cacheable bool) {
	if len(events) == 0 {
		return 0, false
	}
	ttl = finalTTL
	for _, e := range events {
		switch e.Status {
		case StatusInProgress:
			return 0, false
		case StatusFinal:
		default:
			untilKickoff := upcomingTTL
			if !e.Kickoff.IsZero() {
				untilKickoff = max(e.Kickoff.Sub(now), recheckTTL)
			}
			ttl = min(ttl, upcomingTTL, untilKickoff)
		}
	}
	return ttl, true
}

// teamSlug maps an nflverse abbreviation to ESPN's team slug.
func teamSlug(nflverseAbbr string) string {
	switch a := strings.ToUpper(nflverseAbbr); a {
	case "WAS":
		return "wsh"
	case "LA":
		return "lar"
	default:
		return strings.ToLower(a)
	}
}

func (c *Client) fetch(ctx context.Context, url string) ([]Event, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESPN returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Events []rawEvent `json:"events"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(doc.Events))
	for _, re := range doc.Events {
		if e, ok := re.toEvent(); ok {
			events = append(events, e)
		}
	}
	return events, nil
}

type rawEvent struct {
	ID     string    `json:"id"`
	Date   string    `json:"date"`
	Status rawStatus `json:"status"`
	Week   struct {
		Number int `json:"number"`
	} `json:"week"`
	Competitions []struct {
		Status      *rawStatus `json:"status"`
		Venue       rawVenue   `json:"venue"`
		Competitors []struct {
			HomeAway string          `json:"homeAway"`
			Score    json.RawMessage `json:"score"`
			Team     struct {
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
		} `json:"competitors"`
	} `json:"competitions"`
}

type rawStatus struct {
	Type struct {
		State     string `json:"state"`
		Completed bool   `json:"completed"`
	} `json:"type"`
}

type rawVenue struct {
	FullName string `json:"fullName"`
	Address  struct {
		City  string `json:"city"`
		State string `json:"state"`
	} `json:"address"`
	Indoor bool `json:"indoor"`
}

func (re rawEvent) toEvent() (Event, bool) {
	if len(re.Competitions) == 0 {
		return Event{}, false
	}
	comp := re.Competitions[0]

	status := re.Status
	if comp.Status != nil {
		status = *comp.Status
	}
	e := Event{
		ID:     re.ID,
		Week:   re.Week.Number,
		Status: StatusScheduled,
		Venue:  comp.Venue.FullName,
		City:   comp.Venue.Address.City,
		State:  comp.Venue.Address.State,
		Indoor: comp.Venue.Indoor,
	}
	switch {
	case status.Type.Completed:
		e.Status = StatusFinal
	case status.Type.State == "in":
		e.Status = StatusInProgress
	}

	if t, err := time.Parse("2006-01-02T15:04Z", re.Date); err == nil {
		e.Kickoff = t.UTC()
	} else if t, err := time.Parse(time.RFC3339, re.Date); err == nil {
		e.Kickoff = t.UTC()
	}

	for _, cp := range comp.Competitors {
		var score *int
		if e.Status != StatusScheduled {
			score = parseScore(cp.Score)
		}
		switch cp.HomeAway {
		case "home":
			e.HomeAbbr, e.HomeScore = cp.Team.Abbreviation, score
		case "away":
			e.AwayAbbr, e.AwayScore = cp.Team.Abbreviation, score
		}
	}
	return e, true
}

// parseScore accepts both shapes ESPN uses: "23" (scoreboard) and {"value":23.0} (team schedule).
func parseScore(raw json.RawMessage) *int {
	if len(raw) == 0 {
		return nil
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		if n, err := strconv.Atoi(s); err == nil {
			return &n
		}
		return nil
	}
	var obj struct {
		Value *float64 `json:"value"`
	}
	if json.Unmarshal(raw, &obj) == nil && obj.Value != nil {
		n := int(*obj.Value)
		return &n
	}
	return nil
}

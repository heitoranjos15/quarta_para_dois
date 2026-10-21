package nflverse

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/singleflight"

	"quarta-para-dois/backend/internal/cache"
)

type Client struct {
	cache       cache.Cache
	httpClient  *http.Client
	githubToken string
	opts        Options
	sfGroup     singleflight.Group

	games *lru.Cache[string, []Play]

	mu              sync.Mutex
	gameIndex       *GameIndex
	gameIndexAt     time.Time
	seasonFetchedAt map[int]time.Time
}

func NewClient(c cache.Cache, githubToken string, opts Options) *Client {
	opts = opts.withDefaults()
	games, _ := lru.New[string, []Play](opts.MemoryGames) // only errors for size <= 0
	return &Client{
		cache:       c,
		githubToken: githubToken,
		opts:        opts,
		games:       games,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		seasonFetchedAt: map[int]time.Time{},
	}
}

func (c *Client) download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.githubToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.githubToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("rate limited by GitHub")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

package nflverse

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/apache/arrow/go/v16/parquet/file"
	"github.com/apache/arrow/go/v16/parquet/pqarrow"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/singleflight"

	"quarta-para-dois/backend/internal/cache"
)

const (
	schedulesBaseURL = "https://github.com/nflverse/nflverse-data/releases/download/schedules/schedules_%d.parquet"

	cacheTTL = time.Hour
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

func (c *Client) GetSchedules(season int) ([]Game, error) {
	cacheKey := fmt.Sprintf("schedules:%d", season)
	var games []Game
	if c.cacheGet(cacheKey, &games) {
		return games, nil
	}

	result, err, _ := c.sfGroup.Do(cacheKey, func() (any, error) {
		games, err := c.downloadAndParseSchedules(season)
		if err != nil {
			return nil, err
		}
		c.cacheSet(cacheKey, games, cacheTTL)
		return games, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]Game), nil
}

func (c *Client) downloadAndParseSchedules(season int) ([]Game, error) {
	url := fmt.Sprintf(schedulesBaseURL, season)
	return c.downloadAndParseSchedulesURL(url)
}

func (c *Client) downloadAndParseSchedulesURL(url string) ([]Game, error) {
	result, err := c.downloadAndParse(url, func(r *file.Reader) (any, error) {
		return c.parseSchedules(r)
	})
	if err != nil {
		return nil, err
	}
	return result.([]Game), nil
}

func (c *Client) downloadAndParse(url string, parseFn func(*file.Reader) (any, error)) (any, error) {
	data, err := c.download(context.Background(), url)
	if err != nil {
		return nil, err
	}

	pqReader, err := file.NewParquetReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = pqReader.Close() }()

	return parseFn(pqReader)
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

func (c *Client) parseSchedules(r *file.Reader) ([]Game, error) {
	props := pqarrow.ArrowReadProperties{
		BatchSize: 8192,
	}

	arrowReader, err := pqarrow.NewFileReader(r, props, memory.DefaultAllocator)
	if err != nil {
		return nil, err
	}

	var games []Game

	recordReader, err := arrowReader.GetRecordReader(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	defer recordReader.Release()

	for recordReader.Next() {
		record := recordReader.Record()

		numRecords := record.NumRows()
		for i := 0; i < int(numRecords); i++ {
			game := Game{
				GameID:    getString(record, "game_id", i),
				Season:    getInt32(record, "season", i),
				Week:      getInt32(record, "week", i),
				GameType:  getString(record, "game_type", i),
				GameDate:  getString(record, "gameday", i),
				Weekday:   getString(record, "weekday", i),
				GameTime:  getString(record, "gametime", i),
				HomeTeam:  getString(record, "home_team", i),
				AwayTeam:  getString(record, "away_team", i),
				HomeScore: getInt32(record, "home_score", i),
				AwayScore: getInt32(record, "away_score", i),
				Location:  getString(record, "location", i),
				Stadium:   getString(record, "stadium", i),
				Surface:   getString(record, "surface", i),
				Roof:      getString(record, "roof", i),
				Temp:      getInt32(record, "temp", i),
				Wind:      getInt32(record, "wind", i),
				HomeCoach: getString(record, "home_coach", i),
				AwayCoach: getString(record, "away_coach", i),
				Referee:   getString(record, "referee", i),
				StadiumId: getString(record, "stadium_id", i),
			}
			games = append(games, game)
		}
	}

	return games, nil
}

// cacheGet treats a cache error as a miss, after logging it.
func (c *Client) cacheGet(key string, dest any) bool {
	ok, err := c.cache.Get(context.Background(), key, dest)
	if err != nil {
		slog.Warn("cache get", "key", key, "err", err)
		return false
	}
	return ok
}

func (c *Client) cacheSet(key string, val any, ttl time.Duration) {
	if err := c.cache.Set(context.Background(), key, val, ttl); err != nil {
		slog.Warn("cache set", "key", key, "err", err)
	}
}

// Temporary: used by parseSchedules until Task 8 removes it.
func getColumnIndex(record arrow.Record, name string) int {
	schema := record.Schema()
	for i, field := range schema.Fields() {
		if field.Name == name {
			return i
		}
	}
	return -1
}

func getInt32(record arrow.Record, name string, row int) int {
	idx := getColumnIndex(record, name)
	if idx == -1 {
		return 0
	}
	col := record.Column(idx)
	if col == nil {
		return 0
	}

	switch arr := col.(type) {
	case *array.Int32:
		return int(arr.Value(row))
	case *array.Int64:
		return int(arr.Value(row))
	case *array.Float64:
		return int(arr.Value(row))
	case *array.Float32:
		return int(arr.Value(row))
	default:
		return 0
	}
}

func getString(record arrow.Record, name string, row int) string {
	idx := getColumnIndex(record, name)
	if idx == -1 {
		return ""
	}
	col := record.Column(idx)
	if col == nil {
		return ""
	}

	switch arr := col.(type) {
	case *array.String:
		return arr.Value(row)
	case *array.Binary:
		return string(arr.Value(row))
	default:
		return ""
	}
}

# Task 3: NFLverse Parquet Client

**Files:**
- Create: `backend/internal/nflverse/client.go`
- Create: `backend/internal/nflverse/models.go`
- Create: `backend/internal/nflverse/client_test.go`

**Interfaces:**
- Consumes: `Cache`, `config.Config` (GitHubToken)
- Produces: 
  - `GetPBP(season int) ([]Play, error)`
  - `GetSchedules(season int) ([]Game, error)`
  - `GetTeams() ([]Team, error)`
  - `GetRosters(season int) ([]Roster, error)`

---

## Steps

- [ ] **Step 1: Write failing test for PBP download**

```go
// backend/internal/nflverse/client_test.go
func TestClient_GetPBP(t *testing.T) {
    client := nflverse.NewClient(cache, config)
    plays, err := client.GetPBP(2024)
    assert.NoError(t, err)
    assert.Greater(t, len(plays), 0)
    // Verify cache hit on second call
    plays2, _ := client.GetPBP(2024)
    assert.Equal(t, len(plays), len(plays2))
}
```

- [ ] **Step 2: Run test**
Run: `cd backend && go test ./internal/nflverse/... -v -run TestClient_GetPBP`
Expected: FAIL (client not implemented)

- [ ] **Step 3: Implement models.go** — Define `Play`, `Game`, `Team`, `Roster` structs matching nflverse fields (subset of 372 PBP fields needed)

Key fields for `Play`:
```go
type Play struct {
    PlayID            int     `parquet:"play_id"`
    GameID            string  `parquet:"game_id"`
    Quarter           int     `parquet:"qtr"`
    Down              int     `parquet:"down"`
    YardsToGo         int     `parquet:"ydstogo"`
    PlayType          string  `parquet:"play_type"`
    Description       string  `parquet:"desc"`
    EPA               float64 `parquet:"epa"`
    PasserPlayerName  string  `parquet:"passer_player_name"`
    ReceiverPlayerName string  `parquet:"receiver_player_name"`
    RusherPlayerName  string  `parquet:"rusher_player_name"`
    TacklerPlayerName string  `parquet:"tackler_player_name"`
    // ... add fields needed for aggregation
}
```

- [ ] **Step 4: Implement client.go**
- Download: `https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_{year}.parquet`
- Use `github.com/apache/arrow/go/v16/parquet` to read
- Filter by `game_id` when needed
- Cache parsed season data in Redis with key `pbp:{season}`
- Single-flight pattern for concurrent requests (use `sync.Once` or `singleflight`)

- [ ] **Step 5: Run test**
Run: `cd backend && go test ./internal/nflverse/... -v -run TestClient_GetPBP`
Expected: PASS (first run downloads ~180MB, subsequent cached)

- [ ] **Step 6: Repeat for Schedules, Teams, Rosters**

- [ ] **Step 7: Commit**
```bash
git add backend/internal/nflverse/
git commit -m "feat: nflverse parquet client with redis caching"
```
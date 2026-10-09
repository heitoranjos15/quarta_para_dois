# Task 6: HTTP Handlers (Games, Seasons, Teams)

**Files:**
- Create: `backend/internal/handlers/games.go`
- Create: `backend/internal/handlers/seasons.go`
- Create: `backend/internal/handlers/teams.go`
- Create: `backend/internal/handlers/handlers_test.go`
- Modify: `backend/cmd/api/main.go` (register routes)

**Interfaces:**
- Consumes: `NFLVerseClient`, `aggregation` functions
- Produces: HTTP handlers for all endpoints

---

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/seasons` | Available seasons (1999-2024) |
| GET | `/api/seasons/{year}/weeks` | Weeks in season |
| GET | `/api/seasons/{year}/weeks/{week}/games` | Games in week |
| GET | `/api/games/{game_id}` | Game summary |
| GET | `/api/games/{game_id}/stats` | Team comparison stats |
| GET | `/api/games/{game_id}/plays` | Play-by-play (quarter filter) |
| GET | `/api/games/{game_id}/notes` | Game notes, weather, injuries |
| GET | `/api/teams/{abbr}/games?season=2024` | Team's games |

**Response Format:**
```json
{
  "data": { ... },
  "meta": { "cached": true, "season": 2024 }
}
```

---

## Steps

- [ ] **Step 1: Write failing integration test**

```go
// backend/internal/handlers/handlers_test.go
func TestGamesHandler(t *testing.T) {
    router := setupTestRouter() // with mock client
    req := httptest.NewRequest("GET", "/api/games/2024_01_LAC_DEN/stats", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
    var resp APIResponse[TeamComparisonStats]
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NotNil(t, resp.Data)
}
```

- [ ] **Step 2: Run test** → FAIL

- [ ] **Step 3: Implement all handlers**
- Response wrapper: `{ data: T, meta: { cached: bool, season: int } }`
- Proper error handling (404, 503, 429)
- Query params for filtering (quarter, drive)

Example handler structure:
```go
// backend/internal/handlers/games.go
func GameStatsHandler(client *nflverse.Client, agg *aggregation.Engine) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        gameID := chi.URLParam(r, "gameID")
        season := extractSeason(gameID)
        
        plays, cached, err := client.GetPBPForGame(r.Context(), season, gameID)
        if err != nil { ... }
        
        stats := agg.ComputeTeamStats(plays)
        respond(w, stats, cached, season)
    }
}
```

- [ ] **Step 4: Register routes in main.go**

```go
// backend/cmd/api/main.go
r.Route("/api", func(r chi.Router) {
    r.Get("/seasons", seasons.List)
    r.Get("/seasons/{year}/weeks", seasons.Weeks)
    r.Get("/seasons/{year}/weeks/{week}/games", seasons.Games)
    r.Get("/games/{gameID}", games.Detail)
    r.Get("/games/{gameID}/stats", games.Stats)
    r.Get("/games/{gameID}/plays", games.Plays)
    r.Get("/games/{gameID}/notes", games.Notes)
    r.Get("/teams/{abbr}/games", teams.Games)
})
```

- [ ] **Step 5: Run test** → PASS

- [ ] **Step 6: Commit**
```bash
git add backend/internal/handlers/ backend/cmd/api/main.go
git commit -m "feat: http handlers for games, seasons, teams"
```
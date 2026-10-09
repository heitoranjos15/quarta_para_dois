# Quarta Para Dois Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build NFL game analytics portal with Go API proxying nflverse-data Parquet files via GitHub Releases, Redis caching, and React/Vite/TypeScript frontend deployed to GitHub Pages.

**Architecture:** Monorepo with `/backend` (Go 1.23 + Chi + Apache Arrow Parquet + Redis) and `/frontend` (React 18 + TypeScript + Vite + TanStack Query + Tailwind). On-demand data fetching: API downloads season Parquet from nflverse GitHub Releases on first request, caches parsed data in Redis (TTL 1hr), serves JSON to frontend. Frontend routes: Season → Week → Game (tabs: Stats | Play-by-Play | Notes).

**Tech Stack:** Go 1.23, Chi router, github.com/apache/arrow/go/v16/parquet, github.com/redis/go-redis/v9, github.com/kelseyhightower/envconfig, slog; React 18, TypeScript, Vite 5, React Router v6, TanStack Query v5, Zustand, Tailwind CSS, Recharts, ky.

**Spec:** AGENTS.md, CLAUDE.md, README.md

## Global Constraints

- Go 1.23+, React 18, TypeScript 5.6+, Vite 5
- No persistent database — nflverse-data is source of truth
- Redis cache only (TTL 1 hour), in-memory fallback
- Per-season Parquet download on first request (~180MB PBP)
- Monorepo: single repo, separate deployments (API → Fly.io, Frontend → GitHub Pages)
- Shared types: Go models ↔ TypeScript (manual sync)
- Conventional Commits, golangci-lint, ESLint+Prettier
- Performance: cold <10s, cached <200ms, FCP <1.5s, bundle <200KB gz

## Review Focus

1. **GitHub rate limits** — unauthenticated requests limited to 60/hr; must handle 429 with exponential backoff and GITHUB_TOKEN support
2. **Parquet memory** — 2024 PBP ~180MB compressed → ~600MB parsed; must stream/filter not load all if possible, or document memory requirements
3. **Cache stampede** — concurrent requests for same season before cache populated; use single-flight pattern or mutex
4. **Season boundaries** — new season starts September; cache must expire or handle stale data gracefully
5. **Missing data** — nflverse releases may be delayed; API must return 503/404 with clear error, not crash

---

### Task 1: Backend Scaffolding & Config

**Files:**
- Create: `backend/go.mod`
- Create: `backend/go.sum`
- Create: `backend/.env.example`
- Create: `backend/cmd/api/main.go`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/middleware/logging.go`
- Create: `backend/internal/middleware/cors.go`

**Interfaces:**
- Produces: `config.Load() (*Config, error)`, `Config` struct with RedisAddr, GitHubToken, Port, Env

- [ ] **Step 1: Write failing test for config loading**

```go
// backend/internal/config/config_test.go
func TestLoadConfig(t *testing.T) {
    os.Setenv("REDIS_ADDR", "localhost:6379")
    os.Setenv("PORT", "8080")
    cfg, err := config.Load()
    assert.NoError(t, err)
    assert.Equal(t, "localhost:6379", cfg.RedisAddr)
    assert.Equal(t, "8080", cfg.Port)
}
```

- [ ] **Step 2: Run test to verify it fails**
Run: `cd backend && go test ./internal/config/... -v`
Expected: FAIL (package doesn't exist)

- [ ] **Step 3: Implement config in `backend/internal/config/config.go`**

Use `github.com/kelseyhightower/envconfig`. Struct with tags for all env vars.

- [ ] **Step 4: Implement main.go with Chi router, middleware, health endpoint**

```go
// backend/cmd/api/main.go
func main() {
    cfg, _ := config.Load()
    r := chi.NewRouter()
    r.Use(middleware.Logger, middleware.Recoverer, cors.Handler)
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
    http.ListenAndServe(":"+cfg.Port, r)
}
```

- [ ] **Step 5: Run test to verify it passes**
Run: `cd backend && go test ./internal/config/... -v`
Expected: PASS

- [ ] **Step 6: Run main to verify server starts**
Run: `cd backend && go run ./cmd/api`
Expected: "Server starting on :8080"

- [ ] **Step 7: Commit**
```bash
git add backend/
git commit -m "feat: backend scaffolding with config, chi router, middleware"
```

---

### Task 2: Redis Cache Wrapper

**Files:**
- Create: `backend/internal/nflverse/cache.go`
- Create: `backend/internal/nflverse/cache_test.go`

**Interfaces:**
- Consumes: `config.Config` (RedisAddr, RedisPassword)
- Produces: `Cache` interface with `Get(key string) (any, bool)`, `Set(key string, val any, ttl time.Duration)`, `Delete(key string)`

- [ ] **Step 1: Write failing test**

```go
func TestCache_SetGet(t *testing.T) {
    c := cache.NewRedisCache("localhost:6379", "")
    c.Set("test", "value", time.Minute)
    val, ok := c.Get("test")
    assert.True(t, ok)
    assert.Equal(t, "value", val)
}
```

- [ ] **Step 2: Run test (needs Redis)**
Run: `docker run -d -p 6379:6379 redis:7-alpine && cd backend && go test ./internal/nflverse/... -v`
Expected: FAIL (cache not implemented)

- [ ] **Step 3: Implement cache.go**

Use `github.com/redis/go-redis/v9`. JSON marshal/unmarshal for values. Handle connection pooling.

- [ ] **Step 4: Run test**
Run: `cd backend && go test ./internal/nflverse/... -v`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/internal/nflverse/cache.go backend/internal/nflverse/cache_test.go
git commit -m "feat: redis cache wrapper"
```

---

### Task 3: NFLverse Parquet Client

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

- [ ] **Step 1: Write failing test for PBP download**

```go
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

- [ ] **Step 3: Implement models.go** — Define `Play`, `Game`, `Team`, `Roster` structs matching nflverse fields (372 PBP fields → subset needed)

- [ ] **Step 4: Implement client.go**
- Download: `https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_{year}.parquet`
- Use `github.com/apache/arrow/go/v16/parquet` to read
- Filter by `game_id` when needed
- Cache parsed season data in Redis with key `pbp:{season}`
- Single-flight pattern for concurrent requests

- [ ] **Step 5: Run test**
Run: `cd backend && go test ./internal/nflverse/... -v -run TestClient_GetPBP`
Expected: PASS (first run downloads ~180MB, subsequent cached)

- [ ] **Step 6: Repeat for Schedules, Teams, Rosters**

- [ ] **Step 7: Commit**
```bash
git add backend/internal/nflverse/
git commit -m "feat: nflverse parquet client with redis caching"
```

---

### Task 4: Aggregation Layer (Team Stats)

**Files:**
- Create: `backend/internal/aggregation/team_stats.go`
- Create: `backend/internal/aggregation/team_stats_test.go`

**Interfaces:**
- Consumes: `[]Play` (filtered to single game)
- Produces: `TeamComparisonStats` struct with offense/defense metrics

- [ ] **Step 1: Write failing test**

```go
func TestComputeTeamStats(t *testing.T) {
    plays := loadTestPlays("2024_01_LAC_DEN") // fixture
    stats := aggregation.ComputeTeamStats(plays)
    assert.Equal(t, "LAC", stats.HomeTeam)
    assert.Equal(t, "DEN", stats.AwayTeam)
    assert.Greater(t, stats.HomeOffense.EPA_per_play, 0.0)
    assert.Greater(t, stats.AwayDefense.SuccessRate, 0.0)
}
```

- [ ] **Step 2: Run test**
Run: `cd backend && go test ./internal/aggregation/... -v`
Expected: FAIL

- [ ] **Step 3: Implement team_stats.go**
Compute from PBP:
- EPA per play (offense/defense, pass/rush)
- Success rate (epa > 0)
- Red zone efficiency
- 3rd/4th down conversion
- Time of possession
- Explosive plays (20+ pass, 10+ rush)
- Turnover worthy plays

- [ ] **Step 4: Run test**
Run: `cd backend && go test ./internal/aggregation/... -v`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/internal/aggregation/
git commit -m "feat: team comparison stats aggregation"
```

---

### Task 5: Aggregation Layer (Play-by-Play Processing)

**Files:**
- Create: `backend/internal/aggregation/play_stats.go`
- Create: `backend/internal/aggregation/play_stats_test.go`

**Interfaces:**
- Consumes: `[]Play` (single game)
- Produces: `[]PlayDetail` for UI table, grouped by quarter

- [ ] **Step 1: Write failing test**

```go
func TestProcessPlays(t *testing.T) {
    plays := loadTestPlays("2024_01_LAC_DEN")
    details := aggregation.ProcessPlays(plays)
    assert.Len(t, details, 4) // 4 quarters
    assert.Equal(t, 1, details[0].Quarter)
    assert.Contains(t, details[0].Plays[0].Description, "pass")
}
```

- [ ] **Step 2: Run test** → FAIL

- [ ] **Step 3: Implement play_stats.go**
- Filter PBP to game_id
- Group by quarter
- For each play: extract description, EPA, players (passer, receiver, rusher, tackler), play type, down/distance, yardline
- Compute drive info (start, end, result, plays, yards, TOP)

- [ ] **Step 4: Run test** → PASS

- [ ] **Step 5: Commit**

---

### Task 6: HTTP Handlers (Games, Seasons, Teams)

**Files:**
- Create: `backend/internal/handlers/games.go`
- Create: `backend/internal/handlers/seasons.go`
- Create: `backend/internal/handlers/teams.go`
- Create: `backend/internal/handlers/handlers_test.go`
- Modify: `backend/cmd/api/main.go` (register routes)

**Interfaces:**
- Consumes: `NFLVerseClient`, `aggregation` functions
- Produces: HTTP handlers for all endpoints

Endpoints:
- `GET /api/seasons` → []int{1999..2024}
- `GET /api/seasons/{year}/weeks` → []int
- `GET /api/seasons/{year}/weeks/{week}/games` → []Game
- `GET /api/games/{game_id}` → GameDetail
- `GET /api/games/{game_id}/stats` → TeamComparisonStats
- `GET /api/games/{game_id}/plays` → []PlayDetail (with quarter filter)
- `GET /api/games/{game_id}/notes` → GameNotes (weather, injuries, vegas)
- `GET /api/teams/{abbr}/games?season=2024` → []Game

- [ ] **Step 1: Write failing integration test**

```go
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

- [ ] **Step 4: Register routes in main.go**

- [ ] **Step 5: Run test** → PASS

- [ ] **Step 6: Commit**

---

### Task 7: Frontend Scaffolding & API Client

**Files:**
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/api/types.ts`
- Create: `frontend/src/hooks/useGames.ts`
- Create: `frontend/src/hooks/useGameStats.ts`
- Create: `frontend/src/hooks/usePlays.ts`

**Interfaces:**
- Consumes: Vite proxy config (`/api` → `http://localhost:8080`)
- Produces: Typed API functions, React Query hooks

- [ ] **Step 1: Define TypeScript types** matching Go models
- [ ] **Step 2: Implement ky-based API client** with base URL from `import.meta.env.VITE_API_URL`
- [ ] **Step 3: Create React Query hooks** for each endpoint
- [ ] **Step 4: Configure QueryClient** in main.tsx (staleTime 5min, retry 1)
- [ ] **Step 5: Verify types compile** `npm run build`
- [ ] **Step 6: Commit**

---

### Task 8: Season & Week Pages

**Files:**
- Create: `frontend/src/pages/SeasonPage.tsx`
- Create: `frontend/src/pages/WeekPage.tsx`
- Create: `frontend/src/components/game/GameCard.tsx`
- Create: `frontend/src/components/ui/Header.tsx` (already done)

**Interfaces:**
- Consumes: `useGames(season, week)`, `useWeeks(season)`
- Produces: Season grid, Week grid with game cards

- [ ] **Step 1: Implement SeasonPage** — grid of week links
- [ ] **Step 2: Implement WeekPage** — list of GameCards
- [ ] **Step 3: Implement GameCard** — expandable row: "Chargers 33 × 10 Broncos [v]" with team logos, scores, click → `/game/{gameId}`
- [ ] **Step 4: Verify in browser** `npm run dev`
- [ ] **Step 5: Commit**

---

### Task 9: Game Detail Page (Stats Tab)

**Files:**
- Create: `frontend/src/pages/GamePage.tsx`
- Create: `frontend/src/components/stats/TeamComparisonCharts.tsx`
- Create: `frontend/src/components/stats/StatRow.tsx`
- Create: `frontend/src/components/ui/Tabs.tsx`

**Interfaces:**
- Consumes: `useGameStats(gameId)`, `useGameInfo(gameId)`
- Produces: Game header + tabs (Stats | Play-by-Play | Notes)

Stats tab:
- Team logos, names, final score
- Comparison charts (Recharts): EPA/play, Success rate, Pass/Rush split, Red zone, 3rd down
- Stat rows: total yards, TOP, turnovers, penalties

- [ ] **Step 1: Implement GamePage layout** with tabs
- [ ] **Step 2: Implement TeamComparisonCharts** (bar charts, radar chart)
- [ ] **Step 3: Implement StatRow** component
- [ ] **Step 4: Verify with real API** (start backend, visit `/game/2024_01_LAC_DEN`)
- [ ] **Step 5: Commit**

---

### Task 10: Game Detail Page (Play-by-Play Tab)

**Files:**
- Create: `frontend/src/components/pbp/PlayByPlayTable.tsx`
- Create: `frontend/src/components/pbp/QuarterAccordion.tsx`
- Create: `frontend/src/components/pbp/PlayRow.tsx`

**Interfaces:**
- Consumes: `usePlays(gameId, {quarter?, drive?})`
- Produces: Quarter accordions with play rows

Play row: "1 & 10 pass Herbert (QB) to Wilson (TE) 15 yards (EPA: +1.2) — tackled by Hurts (DB)"
- Quarter badge, down/distance, play type icon, description, EPA badge, players involved

- [ ] **Step 1: Implement QuarterAccordion** (collapsible)
- [ ] **Step 2: Implement PlayRow** with formatting
- [ ] **Step 3: Implement PlayByPlayTable** with quarter filter tabs
- [ ] **Step 4: Verify with real data**
- [ ] **Step 5: Commit**

---

### Task 11: Game Detail Page (Notes Tab)

**Files:**
- Create: `frontend/src/components/notes/GameNotes.tsx`
- Create: `frontend/src/components/notes/Weather.tsx`
- Create: `frontend/src/components/notes/Injuries.tsx`
- Create: `frontend/src/components/notes/VegasLines.tsx`

**Interfaces:**
- Consumes: `useGameNotes(gameId)`
- Produces: Game metadata, weather, injuries, vegas lines

- [ ] **Step 1: Implement components**
- [ ] **Step 2: Wire into GamePage Notes tab**
- [ ] **Step 3: Commit**

---

### Task 12: CI/CD & Deployment Config

**Files:**
- Create: `.github/workflows/api.yml`
- Create: `.github/workflows/frontend.yml`
- Create: `.github/workflows/lint.yml`
- Create: `fly.toml` (Fly.io config)

**Interfaces:**
- Produces: Automated deployments on push to main

- [ ] **Step 1: API workflow** — Build Docker, push to Fly.io, set secrets
- [ ] **Step 2: Frontend workflow** — Vite build, deploy to gh-pages
- [ ] **Step 3: Lint workflow** — golangci-lint + ESLint
- [ ] **Step 4: Configure Fly.io app** `flyctl launch --dockerfile backend/Dockerfile`
- [ ] **Step 5: Test deployments**
- [ ] **Step 6: Commit**

---

### Task 13: Docker Compose & Local Dev Polish

**Files:**
- Modify: `docker-compose.yml`
- Create: `backend/Dockerfile`
- Create: `frontend/Dockerfile.dev`
- Create: `Makefile` (optional convenience)

- [ ] **Step 1: Verify `docker-compose up -d` starts API + Redis + Frontend**
- [ ] **Step 2: Test full flow** — browser to localhost:5173 → API → Redis → nflverse
- [ ] **Step 3: Add README.md** with architecture diagram, quick start, deployment guide
- [ ] **Step 4: Commit**

---

### Task 14: Polish & Edge Cases

**Files:**
- Modify: various

- [ ] **Step 1: Error boundaries** in React for failed queries
- [ ] **Step 2: Loading skeletons** for all async components
- [ ] **Step 3: Empty states** (no games, no plays)
- [ ] **Step 4: Responsive design** (mobile week grid, game card)
- [ ] **Step 5: Accessibility** (ARIA labels, keyboard nav, color contrast)
- [ ] **Step 6: Performance** — React.memo, useMemo for charts, virtualized PBP list
- [ ] **Step 7: Run full test suite** `cd backend && go test ./... && cd ../frontend && npm run test && npm run lint`
- [ ] **Step 8: Commit**

---

## Execution Recommendation

**Subagent-driven** recommended because:
- 14 tasks with clear interfaces between them
- Backend and frontend can be developed in parallel after Task 1-3
- Each task produces independently testable code
- Independent review per task catches integration issues early
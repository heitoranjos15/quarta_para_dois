# Part 1: Fix What's Wrong Now — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make every number the app shows correct, stop the Washington/Rams 404s, make production work without paid services, and get CI running so later parts are protected.

**Architecture:** Games come from ESPN's scoreboard, joined to nflverse by ESPN event ID through `games.csv`, which is fetched at runtime and held in memory. Play-by-play is downloaded per season, split per game, kept in a bounded in-memory LRU and in Redis per game. Drive results come from nflverse's own `fixed_drive` / `fixed_drive_result` columns instead of hand-written rules. Upstream clients are injected into handlers through constructors.

**Tech Stack:** Go 1.26, chi v5, Apache Arrow Go v16 (parquet), go-redis v9, testify, miniredis v2 (tests), hashicorp golang-lru v2; React 18, TypeScript, Vite 5, Vitest 2 + Testing Library, ESLint 9; GitHub Actions, Fly.io, GitHub Pages.

**Spec:** `docs/superpowers/specs/2026-10-21-fixes-landing-week-notes-design.md` (Part 1, plus 3.7 freshness rules). Sample upstream payloads: `docs/research/espn-payloads/` (see its README).

## Global Constraints

- Free services only: ESPN public APIs, nflverse (GitHub releases + `nflverse/nfldata`), Open-Meteo, Upstash free tier. No paid APIs.
- No local data files at runtime: `games.csv` is fetched from `https://github.com/nflverse/nfldata/raw/master/data/games.csv` and held in memory; never committed, bundled or read from disk. (Test fixtures under `testdata/` are fine.)
- Canonical game ID everywhere is the nflverse `game_id` (e.g. `2026_06_WAS_SF`). ESPN is called by ESPN event ID, never by IDs rebuilt from abbreviations.
- Cache TTLs: finished games 30 days; upcoming games 1 hour; **games in progress are never cached**; `games.csv` index 24 hours.
- Placeholder data must never be shown as real: missing → omit the field / "Not available".
- Default branch is `master`. Commits are Conventional (`fix:`, `feat:`, `chore:`, `ci:`, `test:`).
- Every bug fix starts with a failing test.
- Backend integration tests that hit the network stay behind `testing.Short()`; CI runs `go test -short ./...`. All new tests use `httptest` servers / miniredis, never the network.

## Review Focus

Inputs no single task's main tests are built around, most likely to bite first. Each line has a pinned test in the owning task.

1. **Game missing from `games.csv`** (just scheduled or flexed): it still lists from ESPN, its page loads, and stats/plays/drives return 404 "not available yet", not 500. Test in Task 8 (`TestService_Week_UnmatchedEventStillListed`).
2. **Redis unreachable in production**: requests still succeed from memory/upstream and a warning is logged; no 500. Test in Task 5 (`TestGamePlays_RedisDown`).
3. **Playoff week URLs** (`/season/2026/week/22`): mapped to ESPN `seasontype=3`, and the Super Bowl is ESPN week **5**, not 4. Test in Task 6 (`TestGameIndex_ESPNWeek`).
4. **A week with a game in progress**: the scoreboard is not cached at all; once every game is final it's cached 30 days. Test in Task 7 (`TestScoreboardTTL`).
5. **Overtime and ties**: an OT drive lands under quarter 5 and ends `END_OF_GAME`; a tied or unplayed game highlights no winner. Tests in Task 9 (`TestComputeDrives_Overtime`) and Task 11 (`GameCard.test.tsx`).

---

### Task 1: Backend baseline — hygiene, logging, server, lint

**Files:**
- Modify: `backend/cmd/api/main.go`, `backend/internal/handlers/handlers.go`, `backend/internal/handlers/handlers_test.go`, `backend/internal/config/config.go`, `backend/Dockerfile`, `backend/go.mod`, `backend/go.sum`
- Delete: `backend/internal/middleware/logging.go`
- Create: `backend/internal/middleware/request_log.go`, `backend/cmd/api/main_test.go`, `backend/.golangci.yml`

**Interfaces:**
- Produces: `func newLogger(env string, w io.Writer) *slog.Logger`, `func newServer(addr string, h http.Handler) *http.Server` (in `cmd/api`); `func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler` (in `middleware`); `handlers.SeasonsGamesHandler(espn espnClientInterface) http.HandlerFunc` (constructor form, like the other handlers).

- [ ] **Step 1: Write failing tests** in `cmd/api/main_test.go`:
  - `TestNewLogger_ProductionIsJSON`: `newLogger("production", &buf).Info("hi")` → `json.Unmarshal(buf.Bytes(), &m)` succeeds and `m["msg"] == "hi"`.
  - `TestNewLogger_DevelopmentIsText`: output for `"development"` is not valid JSON and contains `msg=hi`.
  - `TestNewServer_Timeouts`: `ReadHeaderTimeout == 10*time.Second`, `ReadTimeout == 15*time.Second`, `WriteTimeout == 60*time.Second`, `IdleTimeout == 120*time.Second`, `Addr == ":8080"`.
  - In `handlers_test.go`, change `TestSeasonsGamesHandler` to call `SeasonsGamesHandler(mockESPN)` directly (no context injection) and assert 200 + one game.
- [ ] **Step 2: Run** `cd backend && go test ./cmd/api/ ./internal/handlers/` — expect FAIL (undefined `newLogger`, `newServer`; wrong arity for `SeasonsGamesHandler`).
- [ ] **Step 3: Implement.**
  - `main.go`: build `logger := newLogger(cfg.Env, os.Stdout)` and `slog.SetDefault(logger)`; replace `chimiddleware.Logger` with `middleware.RequestLogger(logger)` (logs method, path, status, duration_ms, request_id); delete the context-injection middleware and the `nflverseClientKey`/`espnClientKey` constants in both `main.go` and `handlers.go`; register `SeasonsGamesHandler(espnClient)`; `if err := newServer(":"+cfg.Port, r).ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { logger.Error(...); os.Exit(1) }`.
  - Remove every `fmt.Printf("DEBUG …")` in `handlers.go` and `nflverse/espn_client.go`; log upstream errors with `slog.Error` before responding 503.
  - `respond`/`respondError`: check the `Encode` error and log it.
  - Delete `middleware/logging.go`, run `go mod tidy` (drops `chi/v4`, `chi v1.5.5`, `x/net/h2c`).
  - `Dockerfile`: delete the `RUN go mod tidy` line.
  - `.golangci.yml`: `version: "2"`, default `standard` linters, `run.tests: true`.
  - `gofmt -w .` for the whole backend.
- [ ] **Step 4: Verify.** `cd backend && go test -short ./... && test -z "$(gofmt -l .)" && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...` → all pass, no lint findings (fix any errcheck/staticcheck findings it reports). `grep -rn "DEBUG\|chi/v4" backend` → no matches.
- [ ] **Step 5: Commit** `chore(api): slog logging, server timeouts, constructor injection, gofmt`.

---

### Task 2: CI, frontend tooling and GitHub Pages

**Files:**
- Modify: `.github/workflows/lint.yml` → rename to `.github/workflows/ci.yml`; `.github/workflows/api.yml`; `.github/workflows/frontend.yml`; `frontend/package.json`; `frontend/vite.config.ts`; `frontend/tsconfig.node.json`; `frontend/src/main.tsx`
- Move: `fly.toml` → `backend/fly.toml` (set `[build] dockerfile = "Dockerfile"`)
- Create: `frontend/eslint.config.js`, `frontend/src/test/setup.ts`, `frontend/src/api/types.test.ts`

**Interfaces:**
- Produces: `npm test` runs Vitest once over `src/**/*.test.{ts,tsx}` with jsdom + jest-dom; `npm run test:e2e` runs Playwright. Later frontend tasks add tests under `src/`.

- [ ] **Step 1: Write a test of existing behavior** `src/api/types.test.ts`: `mapTeamComparisonStats` turns `TimeOfPossession: 125_000_000_000` into `time_of_possession === "2:05"`.
- [ ] **Step 2: Run** `cd frontend && npx vitest run` — expect FAIL: Vitest also collects `tests/*.spec.ts` (Playwright) and errors.
- [ ] **Step 3: Configure.**
  - `vite.config.ts`: `base: process.env.VITE_BASE ?? '/'`; `test: { environment: 'jsdom', include: ['src/**/*.test.{ts,tsx}'], setupFiles: ['src/test/setup.ts'] }` (add `/// <reference types="vitest/config" />`). `setup.ts` imports `@testing-library/jest-dom/vitest`.
  - `package.json`: `"test": "vitest run"`, `"test:watch": "vitest"`, `"test:e2e": "playwright test"`; add devDependency `@types/node`; `tsconfig.node.json` gets `"types": ["node"]`.
  - `main.tsx`: `<BrowserRouter basename={import.meta.env.BASE_URL}>`.
  - `eslint.config.js`: ESLint 9 flat config, i.e. the Vite React-TS template (`@eslint/js` recommended, `typescript-eslint` recommended, `react-hooks`, `react-refresh`), ignoring `dist`, `playwright-report`, `test-results`. Fix whatever it reports in `src/`.
- [ ] **Step 4: Verify locally.** `cd frontend && npm test && npm run lint && npm run build` → all succeed.
- [ ] **Step 5: Workflows.**
  - `ci.yml` (on `push` + `pull_request`): job **backend** (`actions/setup-go@v5` with `go-version-file: backend/go.mod`; `test -z "$(gofmt -l .)"`; `golangci/golangci-lint-action@v8` with `working-directory: backend`, `version: latest`; `go test -short -race ./...`), job **frontend** (Node 20, `npm ci`, `npm run lint`, `npm test`, `npm run build`).
  - `api.yml`: `branches: [master]`; steps `cd backend && flyctl deploy --remote-only`.
  - `frontend.yml`: `branches: [master]`; env `VITE_BASE: /quarta_para_dois/`; after build `cp frontend/dist/index.html frontend/dist/404.html` (deep-link refresh); `peaceiris/actions-gh-pages@v4`.
- [ ] **Step 6: Commit** `ci: run lint, tests and build on master; fix fly and pages deploys`.

---

### Task 3: Shared cache with errors, TLS and context

**Files:**
- Create: `backend/internal/cache/cache.go`, `backend/internal/cache/cache_test.go`
- Delete: `backend/internal/nflverse/cache.go`, `backend/internal/nflverse/cache_test.go`
- Modify: `backend/internal/config/config.go`, `backend/internal/config/config_test.go`, `backend/cmd/api/main.go`, `backend/internal/nflverse/client.go`, `backend/internal/nflverse/espn_client.go` (only to take the shared cache), `README.md` (deploy section), `backend/.env.example`

**Interfaces:**
- Produces:
  ```go
  package cache
  type Cache interface {
      Get(ctx context.Context, key string, dest any) (bool, error) // (false, nil) on miss
      Set(ctx context.Context, key string, val any, ttl time.Duration) error
  }
  func NewRedis(addr, password string, useTLS bool) *Redis // implements Cache, JSON-encoded values
  ```
  `config.Config.RedisTLS bool` (`REDIS_TLS`, default `false`). Callers log a non-nil error with `slog.Warn` and treat it as a miss.

- [ ] **Step 1: Failing tests** (miniredis): `TestRedis_SetGet` round-trips a struct; `TestRedis_Miss` → `(false, nil)`; `TestRedis_TTL` → after `mr.FastForward(ttl+time.Second)` it's a miss; `TestRedis_Down` → after `mr.Close()`, `Get` returns `(false, err != nil)` and `Set` returns an error. `TestLoad_RedisTLS` → `REDIS_TLS=true` sets `RedisTLS`.
- [ ] **Step 2: Run** `go test ./internal/cache/ ./internal/config/` → FAIL (package missing).
- [ ] **Step 3: Implement** `NewRedis` (go-redis `Options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}` when `useTLS`). In `main.go` create one `cache.NewRedis(...)` and pass it to both the nflverse client and the ESPN client; `ESPNClient` stops creating its own `redis.Client` and uses `cache.Cache` (its current cache keys/TTL stay until Task 7).
- [ ] **Step 4: Run** `go test -short ./...` → PASS.
- [ ] **Step 5: Docs.** README "Deploy" section: `fly secrets set REDIS_ADDR=<host>:<port> REDIS_PASSWORD=<token> REDIS_TLS=true` with Upstash free tier values; `FLY_API_TOKEN` repo secret; GitHub Pages source = `gh-pages` branch. Add `REDIS_TLS=false` to `.env.example`.
- [ ] **Step 6: Commit** `fix(api): share one Redis cache, surface cache errors, support TLS`.

---

### Task 4: Play-by-play parquet parsing — projection, nulls, drive columns

**Files:**
- Create: `backend/internal/nflverse/pbp_parse.go`, `backend/internal/nflverse/pbp_parse_test.go`, `backend/internal/nflverse/testutil_test.go` (parquet fixture writer)
- Modify: `backend/internal/nflverse/models.go`, `backend/internal/nflverse/client.go` (remove old `parsePBP`, `getInt32`/`getString`/`getFloat64`, `GetTeams`, `GetRosters`, `Team`, `Roster`; schedules are removed in Task 8)

**Interfaces:**
- Produces:
  ```go
  func ParsePlays(r parquet.ReaderAtSeeker) ([]Play, error)
  ```
  `Play` changes: `EPA *float64` (nil when null); new `FixedDrive int`, `FixedDriveResult string`, `DriveQuarterStart int`, `DriveQuarterEnd int`, `DriveTimeOfPossession string` ("m:ss"), `FieldGoalResult string`, `ExtraPointResult string`, `TwoPointConvResult string`, `TDTeam string`; remove `YardsToGoal`. All other existing fields stay.

- [ ] **Step 1: Fixture helper** `writePlaysParquet(t, rows []map[string]any) []byte` that builds an Arrow table with nflverse's real column types: integer-like columns as **float64** (as in the 2026 file: `qtr`, `down`, `fixed_drive`, …), strings as utf8, all nullable, plus one unused extra column, and writes it with `pqarrow.WriteTable`.
- [ ] **Step 2: Failing tests** in `pbp_parse_test.go`:
  - `TestParsePlays_Fields`: one row with `game_id=2026_06_DAL_GB`, `fixed_drive=5`, `fixed_drive_result="Turnover"`, `drive_quarter_start=1`, `drive_quarter_end=2`, `drive_time_of_possession="2:35"`, `field_goal_result="made"`, `td_team="DAL"`, `epa=0.42` → all mapped (`*p.EPA == 0.42`).
  - `TestParsePlays_NullEPA`: `epa` null → `p.EPA == nil`; null string → `""`; null number → `0`.
  - `TestParsePlays_IntColumnTypes`: the same field read from int32, int64 and float64 columns gives the same `int`.
  - `TestParsePlays_MultipleBatches`: 10 000 rows (more than one 8192 batch) → `len == 10000`, order preserved.
- [ ] **Step 3: Run** `go test ./internal/nflverse/ -run ParsePlays` → FAIL.
- [ ] **Step 4: Implement `ParsePlays`.** Read only the columns `Play` uses (pass their leaf indices to `GetRecordReader`). Resolve each column's index **once per record batch** from the batch schema, check `IsNull(i)` before `Value(i)`, and accept int32/int64/float32/float64 for int fields.
- [ ] **Step 5: Run** the package tests → PASS; `go test -short ./...` still green (fix compile fallout in aggregation from `EPA *float64` with a temporary `epa(p) float64` helper returning 0 for nil; Task 10 replaces it).
- [ ] **Step 6: Commit** `fix(nflverse): read only needed pbp columns, keep nulls, add drive columns`.

---

### Task 5: Per-game play-by-play store and handlers

**Files:**
- Create: `backend/internal/nflverse/pbp_store.go`, `backend/internal/nflverse/pbp_store_test.go`
- Modify: `backend/internal/nflverse/client.go` (`NewClient` signature, delete `GetPBP`), `backend/internal/handlers/handlers.go`, `backend/internal/handlers/handlers_test.go`, `backend/cmd/api/main.go`, `backend/go.mod` (add `github.com/hashicorp/golang-lru/v2`)

**Interfaces:**
- Consumes: `cache.Cache` (Task 3), `ParsePlays` (Task 4).
- Produces:
  ```go
  var ErrGameNotFound = errors.New("game not found")
  type Options struct {
      PBPURLFormat  string        // default "https://github.com/nflverse/nflverse-data/releases/download/pbp/play_by_play_%d.parquet"
      GamesCSVURL   string        // default "https://github.com/nflverse/nfldata/raw/master/data/games.csv" (used in Task 6)
      MemoryGames   int           // default 64
      SeasonRefetch time.Duration // default 1h: min gap between season downloads
      Now           func() time.Time
  }
  func NewClient(c cache.Cache, githubToken string, opts Options) *Client
  func (c *Client) GamePlays(ctx context.Context, gameID string) ([]Play, error)
  ```
  Handlers: `type PlaysSource interface { GamePlays(ctx context.Context, gameID string) ([]nflverse.Play, error) }`; `GameStatsHandler(PlaysSource)`, `GamePlaysHandler(PlaysSource)`, `GameDrivesHandler(PlaysSource)`; `func parseGameID(id string) (season, week int, ok bool)` accepting `^\d{4}_\d{2}_[A-Z]{2,3}_[A-Z]{2,3}$`.

- [ ] **Step 1: Failing store tests** (httptest server serving a fixture parquet with two games, miniredis cache, counting requests):
  - `TestGamePlays_DownloadsOnceAndSplits`: two calls for game A and one for game B → **1** season download; each call returns only that game's plays; Redis has `pbp:game:A` and `pbp:game:B` with TTL 30 days.
  - `TestGamePlays_FromRedis`: new client, same miniredis → returns game A with **0** downloads.
  - `TestGamePlays_UnknownGameThrottled`: unknown game twice within `SeasonRefetch` → `ErrGameNotFound` both times, **1** download; after advancing `Now` past `SeasonRefetch` → a second download.
  - `TestGamePlays_RedisDown` *(Review Focus 2)*: miniredis closed → still returns plays from the download, no error.
  - `TestGamePlays_MemoryBounded`: `MemoryGames: 1`, fetch A then B, close the upstream and miniredis → A now errors (evicted), B is served from memory.
- [ ] **Step 2: Failing handler tests:** `GET /api/games/not-an-id/stats` → **400**; a valid ID whose source returns `ErrGameNotFound` → **404** `{"error":"game not available yet"}`; other errors → 503.
- [ ] **Step 3: Run** `go test ./internal/nflverse/ ./internal/handlers/` → FAIL.
- [ ] **Step 4: Implement `GamePlays`.** Lookup order: memory LRU → Redis `pbp:game:{id}` → season download (singleflight per season, skipped if the last download of that season was less than `SeasonRefetch` ago) → write every game in the file to Redis (30 days) → put the requested game in the LRU → drop the season slice. Season comes from `parseGameID`. Redis errors: `slog.Warn` and continue.
- [ ] **Step 5: Implement handlers.** Validate with `parseGameID` (400), map `ErrGameNotFound` → 404, filter `/plays?drive=` on `FixedDrive`, and delete `filterPlaysByGame` and `extractSeason`. Wire `nflClient` with `nflverse.Options{}` in `main.go`.
- [ ] **Step 6: Run** `go test -short ./...` → PASS.
- [ ] **Step 7: Commit** `perf(api): cache play-by-play per game with bounded memory`.

---

### Task 6: `games.csv` lookup index

**Files:**
- Create: `backend/internal/nflverse/games_index.go`, `backend/internal/nflverse/games_index_test.go`, `backend/internal/nflverse/testdata/games_sample.csv`

**Interfaces:**
- Produces:
  ```go
  type GameRef struct {
      GameID, ESPNID, GameType, HomeTeam, AwayTeam, Roof, Surface, StadiumID, Stadium string
      Season, Week int
  }
  type GameIndex struct{ /* unexported maps */ }
  func ParseGamesCSV(r io.Reader) (*GameIndex, error)          // by header name, not column position
  func (ix *GameIndex) ByGameID(id string) (GameRef, bool)
  func (ix *GameIndex) ByESPNID(id string) (GameRef, bool)
  func (ix *GameIndex) Seasons() []int                         // ascending
  func (ix *GameIndex) Weeks(season int) []int                 // ascending, includes playoff weeks
  func (ix *GameIndex) ESPNWeek(season, week int) (seasonType, espnWeek int, ok bool)
  func (c *Client) GameIndex(ctx context.Context) (*GameIndex, error) // memory only, 24h, singleflight
  ```
  `ESPNWeek`: game type `REG` → `(2, week)`; `WC` → `(3, 1)`; `DIV` → `(3, 2)`; `CON` → `(3, 3)`; `SB` → `(3, 5)`.

- [ ] **Step 1: Fixture** `testdata/games_sample.csv` with the real header (copy the first line of `docs/research/espn-payloads/nflverse/games.csv`, re-downloadable from the URL above) and rows for `2026_06_DAL_GB` (espn `401873007`), `2026_06_WAS_SF` (`401873008`), `2026_06_ARI_LA`, one 2020 week 17 REG game, one 2020 week 18 WC game, one 2026 week 22 SB game, and one row with an empty `espn`.
- [ ] **Step 2: Failing tests:**
  - `TestParseGamesCSV_Lookups`: `ByESPNID("401873008").GameID == "2026_06_WAS_SF"`, `HomeTeam == "SF"`, `Surface`/`Roof` copied; `ByGameID("2026_06_DAL_GB").ESPNID == "401873007"`; empty `espn` → not in the ESPN map.
  - `TestGameIndex_Weeks`: 2020 → includes 17 and 18; 2026 → includes 6 and 22.
  - `TestGameIndex_ESPNWeek` *(Review Focus 3)*: `(2026, 6)` → `(2, 6)`; `(2020, 18)` → `(3, 1)`; `(2026, 22)` → `(3, 5)`; an unknown week → `ok == false`.
  - `TestClient_GameIndexCachedInMemory`: httptest server; two calls → 1 request; after advancing `Now` 25h → 2 requests.
- [ ] **Step 3: Run** → FAIL. **Step 4: Implement** (`encoding/csv`). **Step 5: Run** → PASS.
- [ ] **Step 6: Commit** `feat(nflverse): index games.csv by nflverse and ESPN game id`.

---

### Task 7: ESPN client package

**Files:**
- Create: `backend/internal/espn/client.go`, `backend/internal/espn/client_test.go`, `backend/internal/espn/testdata/scoreboard_week6_final.json`, `backend/internal/espn/testdata/scoreboard_week7_scheduled.json`, `backend/internal/espn/testdata/team_schedule.json`
- Delete: `backend/internal/nflverse/espn_client.go`

**Interfaces:**
- Consumes: `cache.Cache`.
- Produces:
  ```go
  type Status string
  const (StatusScheduled Status = "scheduled"; StatusInProgress Status = "in_progress"; StatusFinal Status = "final")
  type Event struct {
      ID        string
      Kickoff   time.Time
      Status    Status
      HomeAbbr, AwayAbbr string // ESPN abbreviations as returned
      HomeScore, AwayScore *int // nil unless in progress or final
      Venue, City, State string
      Indoor    bool
  }
  type Options struct { BaseURL string; HTTP *http.Client } // default "https://site.api.espn.com/apis/site/v2/sports/football/nfl", 10s timeout
  func NewClient(c cache.Cache, opts Options) *Client
  func (c *Client) Scoreboard(ctx context.Context, season, seasonType, week int) ([]Event, error)
  func (c *Client) TeamSchedule(ctx context.Context, nflverseAbbr string, season int) ([]Event, error)
  func scoreboardTTL(events []Event) (ttl time.Duration, cacheable bool)
  func teamSlug(nflverseAbbr string) string // "WAS"→"wsh", "LA"→"lar", else lowercase
  ```
  Status: `completed == true` → final; `state == "in"` → in progress; else scheduled. Kickoff parsed from the event's `date` (ESPN uses `2026-10-19T00:15Z`). A missing or bad date → zero time, no panic. Cache key `espn:scoreboard:{season}:{seasonType}:{week}`.

- [ ] **Step 1: Fixtures.** Copy `docs/research/espn-payloads/scoreboard-2026-week6.json` and `…week7.json` into `testdata/`, and save one team-schedule response (`…/teams/wsh/schedule?season=2026`).
- [ ] **Step 2: Failing tests:**
  - `TestScoreboard_ParsesWeek6`: 14 events; `401873007` has `HomeAbbr "GB"`, scores 23/13, `StatusFinal`, Kickoff `2026-10-19T00:20Z`, Venue `Lambeau Field`; `401873008` has `AwayAbbr "WSH"`; the request URL has `seasontype=2&week=6&dates=2026`.
  - `TestScoreboard_ScheduledHasNoScores`: week 7 → `HomeScore == nil`, `StatusScheduled`.
  - `TestScoreboardTTL` *(Review Focus 4)*: all final → `(30*24h, true)`; any scheduled and none in progress → `(1h, true)`; any in progress → `(_, false)`. `TestScoreboard_InProgressNotCached`: a fixture with one event `state:"in"` called twice → 2 upstream requests.
  - `TestScoreboard_BadDate`: event `date: ""` → no panic, zero `Kickoff`.
  - `TestTeamSlug`: WAS→wsh, LA→lar, KC→kc.
- [ ] **Step 3: Run** → FAIL. **Step 4: Implement.** **Step 5: Run** → PASS.
- [ ] **Step 6: Commit** `feat(espn): scoreboard and team schedule client with status-aware caching`.

---

### Task 8: Games service and game handlers

**Files:**
- Create: `backend/internal/games/service.go`, `backend/internal/games/service_test.go`
- Modify: `backend/internal/handlers/handlers.go`, `backend/internal/handlers/handlers_test.go`, `backend/cmd/api/main.go`, `backend/internal/nflverse/client.go` and `models.go` (delete `GetSchedules`, `schedulesBaseURL`, `Game`)

**Interfaces:**
- Consumes: `*nflverse.GameIndex` and `GameRef` (Task 6), `espn.Event`/`espn.Status` (Task 7).
- Produces:
  ```go
  package games
  var ErrNotFound = errors.New("game not found")
  type Game struct {
      GameID    string  `json:"game_id"`
      ESPNID    string  `json:"espn_id"`
      Season    int     `json:"season"`
      Week      int     `json:"week"`
      GameType  string  `json:"game_type"`
      Kickoff   string  `json:"kickoff"` // RFC3339 UTC, "" if unknown
      Status    string  `json:"status"`  // scheduled | in_progress | final
      HomeTeam  string  `json:"home_team"` // nflverse abbreviations
      AwayTeam  string  `json:"away_team"`
      HomeScore *int    `json:"home_score"`
      AwayScore *int    `json:"away_score"`
      Stadium   string  `json:"stadium"`
      Location  string  `json:"location"`
      Indoor    bool    `json:"indoor"`
      Roof      string  `json:"roof,omitempty"`
      Surface   string  `json:"surface,omitempty"`
  }
  type IndexSource interface { GameIndex(ctx context.Context) (*nflverse.GameIndex, error) }
  type ESPNSource interface {
      Scoreboard(ctx context.Context, season, seasonType, week int) ([]espn.Event, error)
      TeamSchedule(ctx context.Context, nflverseAbbr string, season int) ([]espn.Event, error)
  }
  func NewService(ix IndexSource, es ESPNSource) *Service
  func (s *Service) Seasons(ctx context.Context) ([]int, error)
  func (s *Service) Weeks(ctx context.Context, season int) ([]int, error)
  func (s *Service) Week(ctx context.Context, season, week int) ([]Game, error)
  func (s *Service) Game(ctx context.Context, gameID string) (Game, error)
  func (s *Service) TeamGames(ctx context.Context, abbr string, season int) ([]Game, error)
  ```
  Join rule: each ESPN event → `ByESPNID(event.ID)`; if found, the ID/teams/week/type/roof/surface come from the `GameRef`. If not found, `HomeTeam`/`AwayTeam` = ESPN abbreviations with `WSH→WAS`, `LAR→LA`, and `GameID = fmt.Sprintf("%d_%02d_%s_%s", season, week, away, home)` (fallback only).
  Handlers: `GameService` interface with the five methods above; `SeasonsListHandler(gs)`, `SeasonsWeeksHandler(gs)`, `SeasonsGamesHandler(gs)`, `GameDetailHandler(gs)`, `GameNotesHandler(gs)`, `TeamGamesHandler(gs)`.

- [ ] **Step 1: Failing service tests** (fake `IndexSource` from `testdata/games_sample.csv` via `ParseGamesCSV`; fake `ESPNSource` from Task 7 fixtures):
  - `TestService_Week_UsesNflverseIDs`: `Week(2026, 6)` contains `GameID "2026_06_WAS_SF"` with `AwayTeam "WAS"`, and `2026_06_ARI_LA` with `HomeTeam "LA"`.
  - `TestService_Week_UnmatchedEventStillListed` *(Review Focus 1)*: an ESPN event with no index row is listed, with `GameID` built from normalized abbreviations.
  - `TestService_Game`: `Game("2026_06_DAL_GB")` calls `Scoreboard(2026, 2, 6)` **once** (no 18-week scan) and returns `Status "final"`, `HomeScore 23`, `Surface` from the index.
  - `TestService_Game_Playoffs`: a SB game ID → `Scoreboard(2026, 3, 5)`.
  - `TestService_TeamGames_IDs`: no ID contains `_00_`; all IDs come from the index.
  - `TestService_Weeks_PerSeason`: 2020 → includes 17 and 18; empty for an unknown season.
- [ ] **Step 2: Failing handler tests:** `/api/seasons/2026/weeks/6/games` → 200 with `data[].game_id`; `/api/games/2026_06_XXX_YYY` with `ErrNotFound` → 404; `/api/games/bad` → 400; `/api/games/{id}/notes` response has no `vegas` and no `condition` keys (no placeholders), only `stadium`, `roof`, `surface`, `indoor`.
- [ ] **Step 3: Run** → FAIL. **Step 4: Implement** the service and handlers; wire `games.NewService(nflClient, espnClient)` in `main.go`; delete the dead schedule code. **Step 5: Run** `go test -short ./...` → PASS.
- [ ] **Step 6: Manual check.** `cd backend && go run ./cmd/api` with local Redis (`docker compose up -d redis`), then `curl -s localhost:8080/api/seasons/2026/weeks/6/games | jq '.data[] | .game_id'` → includes `2026_06_WAS_SF`, and `curl -s -o /dev/null -w '%{http_code}' localhost:8080/api/games/2026_06_WAS_SF/stats` → `200`.
- [ ] **Step 7: Commit** `fix(api): join ESPN games to nflverse by ESPN id; fixes WAS/LA game pages`.

---

### Task 9: Drive aggregation from nflverse drive columns

**Files:**
- Modify: `backend/internal/aggregation/drive_stats.go`, `backend/internal/aggregation/drive_stats_test.go`, `backend/internal/aggregation/play_stats.go`, `backend/internal/aggregation/play_stats_test.go`, `backend/internal/handlers/handlers.go` (drop the duplicate `DriveSummaryResponse`; keep the `aggregation` one, drop its unused `Overtime`)

**Interfaces:**
- Consumes: `Play.FixedDrive`, `FixedDriveResult`, `DriveQuarterStart`, `DriveQuarterEnd`, `DriveTimeOfPossession`, `FieldGoalResult`, `ExtraPointResult`, `TwoPointConvResult` (Task 4).
- Produces: `ComputeDrives(plays []nflverse.Play, homeTeam, awayTeam string) map[int]QuarterDriveSummary` (same signature); `func driveResult(fixedDriveResult string, plays []nflverse.Play) DriveEndReason`; `func drivePoints(result DriveEndReason, plays []nflverse.Play) int`. `DriveEndReason` values: `TOUCHDOWN, FIELD_GOAL, MISSED_FIELD_GOAL, PUNT, INTERCEPTION, FUMBLE_LOST, DOWNS, SAFETY, END_OF_HALF, END_OF_GAME, OTHER` (`EXTRA_POINT` removed).

  Mapping from `fixed_drive_result`: `Touchdown`→TOUCHDOWN; `Field goal`→FIELD_GOAL; `Missed field goal`→MISSED_FIELD_GOAL; `Punt`→PUNT; `Turnover on downs`→DOWNS; `Safety`→SAFETY; `Turnover` and `Opp touchdown` → INTERCEPTION if any play has `Interception==1`, else FUMBLE_LOST if any `FumbleLost==1`, else PUNT if the last play is a punt, else OTHER; `End of half` → END_OF_GAME when `DriveQuarterEnd >= 4`, else END_OF_HALF; anything else → OTHER.
  Points (offense only): TOUCHDOWN = 6 + 1 if any `ExtraPointResult=="good"` + 2 if any `TwoPointConvResult=="success"`; FIELD_GOAL = 3; everything else 0.
  Grouping: by `FixedDrive` across the whole game. The drive is listed under quarter `DriveQuarterStart` (fallback: first play's `Quarter`) and counted once in that quarter's totals. Team = first non-empty `Posteam`. TOP = `DriveTimeOfPossession`.

- [ ] **Step 1: Fix the fixture and write failing tests:**
  - `TestComputeDrives_FourthDownPuntIsPunt`: last play `Down 4, PlayType "punt"`, `fixed_drive_result "Punt"` → PUNT, 0 points.
  - `TestComputeDrives_FieldGoalMadeAndMissed`: "Field goal" → 3 points; "Missed field goal" → MISSED_FIELD_GOAL, 0 points.
  - `TestComputeDrives_PickSix`: "Opp touchdown" with `Interception 1`, `TDTeam` = the defense → INTERCEPTION, 0 points for the offense.
  - `TestComputeDrives_TouchdownPoints`: XP good → 7; two-point success → 8; XP failed → 6.
  - `TestComputeDrives_CrossQuarter`: plays in Q1 and Q2 with one `FixedDrive 5`, `DriveQuarterStart 1` → exactly one drive, under quarter 1, `PlaysCount` = all its plays; quarter 2 totals don't count it.
  - `TestComputeDrives_Overtime` *(Review Focus 5)*: `DriveQuarterStart 5`, "End of half", `DriveQuarterEnd 5` → under key `5`, END_OF_GAME.
  - `TestProcessPlays_UsesSameDriveResult`: `ProcessPlays` drive summary for a 4th-down punt → `"PUNT"`.
- [ ] **Step 2: Run** `go test ./internal/aggregation/` → FAIL.
- [ ] **Step 3: Implement**; `play_stats.go` reuses `driveResult` (delete its own rule chain and the `YardsGained > 0` FG guess).
- [ ] **Step 4: Run** `go test -short ./...` → PASS.
- [ ] **Step 5: Commit** `fix(aggregation): use nflverse drive results; correct punts, field goals and points`.

---

### Task 10: Team stats on scrimmage plays only

**Files:**
- Modify: `backend/internal/aggregation/team_stats.go`, `backend/internal/aggregation/team_stats_test.go`

**Interfaces:**
- Consumes: `Play.EPA *float64`, `FixedDrive`, `FixedDriveResult`, `DriveTimeOfPossession`.
- Produces: `TeamUnitStats` gains `RedZoneTrips int`, `RedZoneTDs int`; `RedZonePct = RedZoneTDs / RedZoneTrips` (0 when no trips); `TimeOfPossession` (still `time.Duration`) = the sum of `DriveTimeOfPossession` over the unit's distinct `FixedDrive`s. `func isScrimmage(p nflverse.Play) bool` = `PlayType == "pass" || PlayType == "run"`. Remove the temporary `epa()` helper from Task 4.

- [ ] **Step 1: Failing tests:**
  - `TestComputeTeamStats_IgnoresNonPlays`: 4 scrimmage plays with EPA `+1, -1, +1, +1` plus a timeout row, a kickoff, an extra point and a `no_play` penalty (EPA `+5` each) → `Plays == 4`, `EPA_per_play == 0.5`, `SuccessRate == 0.75`.
  - `TestComputeTeamStats_NilEPA`: one scrimmage play with nil EPA → excluded from the EPA average but counted in `Plays` and `PassRate`.
  - `TestComputeTeamStats_RedZoneTrips`: two drives reach `Yardline100 <= 20` (one "Touchdown", one "Field goal"), and one never does → `RedZoneTrips 2`, `RedZoneTDs 1`, `RedZonePct 0.5`.
  - `TestComputeTeamStats_TOP`: drives `"6:30"` and `"2:07"` → `TimeOfPossession == 8*time.Minute + 37*time.Second`.
- [ ] **Step 2: Run** → FAIL. **Step 3: Implement.** **Step 4: Run** `go test -short ./...` → PASS.
- [ ] **Step 5: Commit** `fix(aggregation): compute team rates on scrimmage plays; red zone per trip`.

---

### Task 11: Frontend fixes

**Files:**
- Modify: `frontend/src/api/types.ts`, `frontend/src/api/client.ts`, `frontend/src/components/game/GameCard.tsx`, `frontend/src/components/ui/TeamLogo.tsx`, `frontend/src/components/stats/StatRow.tsx`, `frontend/src/pages/GamePage.tsx`, `frontend/src/components/pbp/PlayRow.tsx` (null EPA), `frontend/src/components/pbp/DriveItem.tsx` (result colors for the new values)
- Create: `frontend/src/components/game/GameCard.test.tsx`, `frontend/src/components/ui/TeamLogo.test.tsx`, `frontend/src/components/stats/StatRow.test.tsx`

**Interfaces:**
- Consumes: the `games.Game` JSON (Task 8), `TeamUnitStats` additions (Task 10), `PlayDetail.EPA` nullable, the new `DriveEndReason` strings (Task 9).
- Produces:
  - `Game` = the Task 8 JSON fields as TS (`home_score: number | null`, `status: 'scheduled' | 'in_progress' | 'final'`, …). `GameAPI` and `mapGameAPItoGame` are deleted because the backend now sends snake_case.
  - `TeamUnitStats` gains `red_zone_trips`, `red_zone_tds`, `time_of_possession_seconds: number`.
  - `PlayDetail.epa: number | null`.
  - `StatRow` gains `format?: (n: number) => string`.

- [ ] **Step 1: Failing tests:**
  - `TeamLogo.test.tsx`: renders the full abbreviation, so `NYG`, `NYJ`, `LAC`, `LA` and `LV` render five different texts; `LA` and `WAS` get team colors (not the grey fallback `#6B7280`).
  - `GameCard.test.tsx` *(Review Focus 5)*: final 23–13 → only the home side has `team-score-winner`; final 20–20 → neither; `status: 'scheduled'` with null scores → neither, and scores render as `–`.
  - `StatRow.test.tsx`: `home=517, away=383, format=m:ss` → renders `8:37` and `6:23`, with the home value highlighted as better.
- [ ] **Step 2: Run** `cd frontend && npm test` → FAIL.
- [ ] **Step 3: Implement.**
  - Update `Game` consumers (`WeekPage`, `GameCard`, `GamePage`): `venue` → `stadium`, `game_date` → `kickoff` (format with `toLocaleString`, and show nothing when it's `""`). Update `GameNotesAPI` to `{ stadium, roof, surface, indoor }`.
  - `GamePage`: replace the four one-sided `StatSection`s with an **Offense** section (home offense vs away offense) and a **Defense** section (home defense vs away defense), both values filled. Time of possession uses `time_of_possession_seconds` with the `m:ss` format. The label "Red Zone %" becomes "Red Zone TD %" with `tds/trips` shown.
  - `GameNotesPanel`: render only the fields the API returns (stadium, roof, surface, indoor), with no attendance, humidity, injuries, officials or weather placeholders. Part 3 replaces this panel.
  - `PlayRow`: show `—` for null EPA.
  - `DriveItem`: colors for `MISSED_FIELD_GOAL`, `INTERCEPTION`, `FUMBLE_LOST`, `DOWNS`; drop `EXTRA_POINT`.
- [ ] **Step 4: Verify.** `npm test && npm run lint && npm run build` → pass. Then `npm run dev` with the API running, and open `http://localhost:5173/season/2026/week/6`: the WAS @ SF card opens a game page with stats, and DAL @ GB's drive list shows punts as PUNT and field goals with points.
- [ ] **Step 5: Commit** `fix(frontend): two-sided stats, unique team badges, winner only for decided games`.

---

## Out of scope for this plan

- `GET /api/current-week`, landing-week redirect, playoff week labels in the UI: Part 2.
- Injury, weather, referee and odds endpoints and the new Notes tab: Part 3. Task 8 only removes the placeholders from `/notes`.
- Line-movement history: issue #7.

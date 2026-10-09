# AGENTS.md - Development Guidelines for quarta_para_dois

## Project Overview
NFL game analytics portal - React frontend + Go API proxying nflverse-data (Parquet on GitHub Releases) with Redis caching.

## Tech Stack

### Backend (Go)
- **Language**: Go 1.23+
- **Framework**: Chi router (lightweight, stdlib-compatible)
- **Parquet**: Apache Arrow Go (github.com/apache/arrow/go/v16/parquet)
- **Cache**: Redis (github.com/redis/go-redis/v9)
- **HTTP Client**: stdlib net/http with custom transport
- **Config**: envconfig (github.com/kelseyhightower/envconfig)
- **Logging**: slog (stdlib)
- **Testing**: testify (github.com/stretchr/testify)

### Frontend (React + Vite + TypeScript)
- **Framework**: React 18 + TypeScript
- **Build**: Vite 5
- **Routing**: React Router v6
- **State**: TanStack Query v5 (server state) + Zustand (client state)
- **UI**: Tailwind CSS + Headless UI / Radix UI
- **Charts**: Recharts or Apache ECharts
- **HTTP**: Ky or fetch with TanStack Query
- **Testing**: Vitest + React Testing Library

### Infrastructure
- **API Hosting**: Fly.io (Docker)
- **Frontend Hosting**: GitHub Pages (via Actions)
- **Cache**: Redis (Upstash or Fly.io Redis)
- **CI/CD**: GitHub Actions
- **Source of Truth**: nflverse/nflverse-data GitHub Releases

## Code Structure

```
quarta_para_dois/
├── backend/                    # Go module
│   ├── cmd/
│   │   └── api/               # Main entry point
│   ├── internal/
│   │   ├── nflverse/          # nflverse data client
│   │   │   ├── client.go      # Download & parse Parquet
│   │   │   ├── models.go      # Domain models
│   │   │   └── cache.go       # Redis cache wrapper
│   │   ├── handlers/          # HTTP handlers
│   │   │   ├── games.go       # Game endpoints
│   │   │   ├── seasons.go     # Season/week endpoints
│   │   │   └── teams.go       # Team endpoints
│   │   ├── aggregation/       # Stats computation
│   │   │   ├── team_stats.go  # Team comparison stats
│   │   │   └── play_stats.go  # Play-by-play processing
│   │   └── middleware/        # HTTP middleware
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .env.example
├── frontend/                   # React + Vite + TS
│   ├── src/
│   │   ├── components/        # Reusable UI components
│   │   │   ├── game/          # GameCard, GameHeader, etc.
│   │   │   ├── stats/         # TeamComparisonCharts, StatRow
│   │   │   ├── pbp/           # PlayByPlayTable, QuarterAccordion
│   │   │   ├── notes/         # GameNotes, Weather, Injuries
│   │   │   └── ui/            # Button, Card, Tabs, Loading
│   │   ├── pages/             # Route-level pages
│   │   │   ├── SeasonPage.tsx
│   │   │   ├── WeekPage.tsx
│   │   │   └── GamePage.tsx
│   │   ├── hooks/             # Custom React hooks
│   │   │   ├── useGames.ts
│   │   │   ├── useGameStats.ts
│   │   │   └── usePlays.ts
│   │   ├── api/               # API client + types
│   │   │   ├── client.ts
│   │   │   └── types.ts
│   │   ├── types/             # Shared TypeScript types
│   │   ├── utils/             # Helpers
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── index.css
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── index.html
├── docker-compose.yml          # Local dev (API + Redis)
├── .github/workflows/          # CI/CD
│   ├── api.yml                # Build & deploy API
│   ├── frontend.yml           # Build & deploy to Pages
│   └── lint.yml               # Lint both
├── README.md
├── AGENTS.md
├── CLAUDE.md
└── docs/
    └── superpowers/specs/     # Design specs
```

## Data Pipeline Flow

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │────▶│  React (Vite)    │────▶│  Go API     │────▶│  Redis      │
│   Browser   │     │  (GitHub Pages)  │     │  (Fly.io)   │     │  (Upstash)  │
└─────────────┘     └──────────────────┘     └──────┬──────┘     └──────┬──────┘
                                                    │                 │
                                                    ▼                 │
                                           ┌──────────────────┐       │
                                           │ nflverse-data    │◀──────┘
                                           │ GitHub Releases  │
                                           │ (Parquet files)  │
                                           └──────────────────┘
```

**Request Flow:**
1. React calls `GET /api/games/{gameId}/stats`
2. Go API checks Redis for `pbp:2024` (season cache)
3. Cache miss → Download `play_by_play_2024.parquet` from GitHub Releases
4. Parse with Apache Arrow, filter to `game_id`
5. Compute team comparison stats (offense/defense EPA, success rate, etc.)
6. Store parsed season data in Redis (TTL: 1 hour)
7. Return JSON to React
8. Subsequent requests for same season hit Redis cache

## Key Design Decisions

1. **No persistent database** - nflverse-data is source of truth
2. **On-demand fetching** - Download season data only when first requested
3. **Redis as ephemeral cache** - TTL-based, auto-expires
4. **Per-season caching** - First request for season downloads ~180MB Parquet
5. **Monorepo** - Single repo, separate deployments
6. **Shared types** - Go models ↔ TypeScript types (manual sync or codegen)

## Development Commands

```bash
# Backend
cd backend
go mod tidy
go run ./cmd/api           # Run locally (needs Redis)
go test ./...
go vet ./...
golangci-lint run

# Frontend
cd frontend
npm install
npm run dev                # Vite dev server
npm run build              # Production build
npm run test               # Vitest
npm run lint               # ESLint

# Docker (local full stack)
docker-compose up -d       # Starts API + Redis
```

## Environment Variables

### Backend (.env)
```
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
GITHUB_TOKEN=ghp_xxx       # Optional: higher rate limits
PORT=8080
ENV=development
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
```

## API Contract

See `backend/internal/handlers/` for endpoint definitions.
All responses follow: `{ "data": T, "meta": { "cached": bool, "season": int } }`

## Deployment

### API (Fly.io)
```bash
flyctl launch --name quarta-para-dois-api --dockerfile backend/Dockerfile
flyctl secrets set REDIS_ADDR=... REDIS_PASSWORD=... GITHUB_TOKEN=...
flyctl deploy
```

### Frontend (GitHub Pages)
- Push to `main` triggers `.github/workflows/frontend.yml`
- Builds Vite app → deploys to `gh-pages` branch
- Available at `https://heitor.github.io/quarta_para_dois/`

## Code Style

- **Go**: Standard `gofmt`, `golangci-lint` (default config)
- **TypeScript**: ESLint + Prettier (Airbnb base)
- **Commits**: Conventional Commits (`feat:`, `fix:`, `chore:`)
- **PRs**: Require CI pass + 1 review

## Testing Strategy

- **Backend**: Unit tests for aggregation logic, integration tests for handlers
- **Frontend**: Component tests (React Testing Library), hook tests (Vitest)
- **E2E**: Playwright (optional, for critical flows)

## Performance Targets

- First request (cold cache): < 10s (download + parse)
- Cached requests: < 200ms
- Frontend FCP: < 1.5s
- Bundle size: < 200KB gzipped
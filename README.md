# Quarta Para Dois - NFL Analytics Portal

> **Quarta Para Dois** (Fourth Down for Two) - A modern NFL game analytics portal built with Go and React.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER BROWSER                                    │
│                         (GitHub Pages - Static)                             │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │ HTTPS /api/*
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GO API                                          │
│                        (Fly.io - Docker Container)                          │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────────────────────┐  │
│  │  Chi Router │──│  Handlers    │──│  Aggregation Engine              │  │
│  │  + Middleware     (Games,      │  │  • Team Comparison Stats       │  │
│  │               Seasons,        │  │  • Play-by-Play Processing     │  │
│  │               Teams)          │  │  • Drive Analysis              │  │
│  └─────────────┘  └──────────────┘  └──────────────────────────────────┘  │
│         │                │                    │                            │
│         │                │                    ▼                            │
│         │                │  ┌─────────────────────────────────────────┐  │
│         │                │  │       NFLverse Client                   │  │
│         │                │  │  • Download Parquet from GitHub         │  │
│         │                │  │  • Parse with Apache Arrow              │  │
│         │                │  │  • Filter by game_id                    │  │
│         │                │  └─────────────────────────────────────────┘  │
│         │                │                    │                            │
│         ▼                ▼                    ▼                            │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                      REDIS CACHE (Upstash/Fly.io)                    │  │
│  │  Keys: pbp:{season}, sched:{season}, teams, rosters:{season}        │  │
│  │  TTL: 1 hour (auto-expires)                                         │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────┬───────────────────────────────────────────┘
                                  │ HTTPS (GitHub Releases API)
                                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        NFLVERSE-DATA (Source of Truth)                      │
│                    https://github.com/nflverse/nflverse-data                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │ PBP         │  │ Schedules   │  │ Teams       │  │ Rosters     │ ...   │
│  │ play_by_    │  │ sched_      │  │ teams.      │  │ weekly_     │       │
│  │ play_2024.  │  │ 2024.       │  │ parquet     │  │ rosters_    │       │
│  │ parquet     │  │ parquet     │  │             │  │ 2024.parquet│       │
│  │ (~180MB)    │  │ (~2MB)      │  │ (~1MB)      │  │ (~15MB)     │       │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Data Pipeline Flow

```
Request: GET /api/games/2024_01_LAC_DEN/stats
                    │
                    ▼
         ┌─────────────────────┐
         │ Check Redis Cache   │──Hit──▶ Return cached JSON
         │ Key: pbp:2024       │
         └──────────┬──────────┘
                    │ Miss
                    ▼
         ┌─────────────────────┐
         │ Download Parquet    │
         │ github.com/nflverse/│
         │ nflverse-data/      │
         │ releases/download/  │
         │ pbp/play_by_play_   │
         │ 2024.parquet        │
         └──────────┬──────────┘
                    │ (~180MB, ~3-5s)
                    ▼
         ┌─────────────────────┐
         │ Parse with Apache   │
         │ Arrow (Go)          │
         │ Filter game_id =    │
         │ "2024_01_LAC_DEN"   │
         └──────────┬──────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │ Compute Team Stats  │
         │ • EPA/play (off/def)│
         │ • Success Rate      │
         │ • Red Zone %        │
         │ • 3rd/4th Down %    │
         │ • TOP, Explosives   │
         └──────────┬──────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │ Store in Redis      │
         │ Key: pbp:2024       │
         │ TTL: 1 hour         │
         └──────────┬──────────┘
                    │
                    ▼
         ┌─────────────────────┐
         │ Return JSON         │
         │ { data: {...},      │
         │   meta: {cached:    │
         │   false, season:    │
         │   2024} }           │
         └─────────────────────┘
```

## Tech Stack

### Backend (Go)
| Component | Technology |
|-----------|------------|
| Language | Go 1.23+ |
| Router | Chi (github.com/go-chi/chi/v5) |
| Parquet | Apache Arrow Go (github.com/apache/arrow/go/v16/parquet) |
| Cache | Redis (github.com/redis/go-redis/v9) |
| Config | envconfig (github.com/kelseyhightower/envconfig) |
| Logging | slog (stdlib) |
| Testing | testify (github.com/stretchr/testify) |
| Linting | golangci-lint |

### Frontend (React + TypeScript)
| Component | Technology |
|-----------|------------|
| Framework | React 18 + TypeScript |
| Build | Vite 5 |
| Routing | React Router v6 |
| Server State | TanStack Query v5 |
| Client State | Zustand |
| Styling | Tailwind CSS |
| Charts | Recharts |
| HTTP | Ky |
| Testing | Vitest + React Testing Library |
| Linting | ESLint + Prettier (Airbnb) |

### Infrastructure
| Component | Platform |
|-----------|----------|
| API Hosting | Fly.io (Docker) |
| Frontend Hosting | GitHub Pages (GitHub Actions) |
| Cache | Redis (Upstash or Fly.io Redis) |
| CI/CD | GitHub Actions |
| Data Source | nflverse/nflverse-data GitHub Releases |

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
│   │   │   ├── game/          # GameCard, GameHeader
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
    └── superpowers/
        ├── specs/             # Design specs
        └── plans/             # Implementation plans
```

## API Endpoints

| Endpoint | Description | Cache Key |
|----------|-------------|-----------|
| `GET /api/seasons` | Available seasons (1999-2024) | — |
| `GET /api/seasons/{year}/weeks` | Weeks in season | `sched:{year}` |
| `GET /api/seasons/{year}/weeks/{week}/games` | Games in week | `sched:{year}` |
| `GET /api/games/{game_id}` | Game summary | `sched:{year}` |
| `GET /api/games/{game_id}/stats` | Team comparison stats | `pbp:{year}` |
| `GET /api/games/{game_id}/plays` | Play-by-play (quarter filter) | `pbp:{year}` |
| `GET /api/games/{game_id}/notes` | Game notes, weather, injuries | `sched:{year}` + `pbp:{year}` |
| `GET /api/teams/{abbr}/games?season=2024` | Team's games | `sched:{year}` |

**Response Format:**
```json
{
  "data": { ... },
  "meta": {
    "cached": true,
    "season": 2024
  }
}
```

## Development

### Prerequisites
- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- Redis (via docker-compose)

### Quick Start

```bash
# 1. Clone and enter
cd quarta_para_dois

# 2. Start local stack (API + Redis)
docker-compose up -d

# 3. Frontend dev server (separate terminal)
cd frontend && npm install && npm run dev
# → http://localhost:5173

# 4. Backend API
cd backend && go mod tidy && go run ./cmd/api
# → http://localhost:8080
```

### Environment Variables

**Backend** (`backend/.env`):
```env
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
GITHUB_TOKEN=ghp_xxx    # Optional: 5000 req/hr vs 60
PORT=8080
ENV=development
```

**Frontend** (`frontend/.env`):
```env
VITE_API_URL=http://localhost:8080
```

### Testing

```bash
# Backend
cd backend
go test ./...
go vet ./...
golangci-lint run

# Frontend
cd frontend
npm run test
npm run lint
npm run build
```

## Deployment

### API (Fly.io)
```bash
# One-time setup
flyctl launch --name quarta-para-dois-api --dockerfile backend/Dockerfile
flyctl secrets set REDIS_ADDR=... REDIS_PASSWORD=... GITHUB_TOKEN=...

# Deploy (automatic via GitHub Actions on push to main)
flyctl deploy
```

### Frontend (GitHub Pages)
- Push to `main` triggers `.github/workflows/frontend.yml`
- Builds Vite app → deploys to `gh-pages` branch
- Available at `https://<username>.github.io/quarta_para_dois/`

## Data Sources

All data from **nflverse/nflverse-data** GitHub Releases:

| Dataset | Release Tag | Files | Update Frequency |
|---------|-------------|-------|------------------|
| Play-by-Play | `pbp` | `play_by_play_{year}.parquet` | Weekly (in-season) |
| Schedules | `schedules` | `sched_{year}.parquet` | Pre-season + weekly |
| Teams | `teams` | `teams.parquet` | Rarely |
| Rosters | `weekly_rosters` | `weekly_rosters_{year}.parquet` | Weekly |
| Team Stats | `stats_team` | `team_stats_{type}_{year}.parquet` | Weekly |
| Player Stats | `stats_player` | `player_stats_{type}_{year}.parquet` | Weekly |

## Performance Targets

| Metric | Target |
|--------|--------|
| First request (cold cache) | < 10s |
| Cached requests | < 200ms |
| Frontend FCP | < 1.5s |
| Bundle size (gz) | < 200KB |
| Memory (per season) | ~700MB |

## License

CC-BY-4.0 — Data from nflverse. Code: MIT.
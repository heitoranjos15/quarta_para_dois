# CLAUDE.md - Claude Code Instructions for quarta_para_dois

## Project Overview
NFL game analytics portal with Go backend + React frontend. Monorepo architecture.

## Quick Start

```bash
# Backend
cd backend && go mod tidy && go run ./cmd/api

# Frontend
cd frontend && npm install && npm run dev

# Full stack (requires Docker)
docker-compose up -d
```

## Architecture Summary

- **Backend**: Go 1.23 + Chi router + Apache Arrow (Parquet) + Redis cache
- **Frontend**: React 18 + TypeScript + Vite + TanStack Query + Tailwind
- **Data Source**: nflverse/nflverse-data GitHub Releases (Parquet)
- **Deployment**: API → Fly.io, Frontend → GitHub Pages

## Key Files

| Purpose | File |
|---------|------|
| API Entry | `backend/cmd/api/main.go` |
| NFLverse Client | `backend/internal/nflverse/client.go` |
| Redis Cache | `backend/internal/nflverse/cache.go` |
| Game Handlers | `backend/internal/handlers/games.go` |
| Aggregation | `backend/internal/aggregation/team_stats.go` |
| Frontend Entry | `frontend/src/main.tsx` |
| API Client | `frontend/src/api/client.ts` |
| Game Page | `frontend/src/pages/GamePage.tsx` |

## Common Tasks

### Add New API Endpoint
1. Define handler in `backend/internal/handlers/`
2. Register route in `backend/cmd/api/main.go`
3. Add TypeScript types in `frontend/src/api/types.ts`
4. Create React Query hook in `frontend/src/hooks/`

### Add New Aggregation
1. Create function in `backend/internal/aggregation/`
2. Call from handler
3. Add response types

### Modify UI
1. Edit components in `frontend/src/components/`
2. Pages in `frontend/src/pages/`
3. Styles via Tailwind classes

## Testing

```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && npm run test
```

## Linting

```bash
# Backend
cd backend && golangci-lint run

# Frontend
cd frontend && npm run lint
```

## Environment

Copy `.env.example` to `.env` in both `backend/` and `frontend/`.

## Git Workflow

- Branch: `feat/`, `fix/`, `chore/`
- Commits: Conventional (`feat: add game stats endpoint`)
- PR: Requires CI pass + review

## Debugging

- API logs: `slog` structured logging (JSON in prod)
- Frontend: React DevTools + TanStack Query DevTools
- Redis: `redis-cli MONITOR` or Upstash dashboard
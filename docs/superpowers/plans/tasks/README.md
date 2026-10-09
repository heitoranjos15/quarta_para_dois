# Task Index - Quarta Para Dois Implementation

Each task is in a separate file for independent session work.

## Backend Tasks

| Task | File | Description |
|------|------|-------------|
| 1 | [task-01-backend-scaffolding.md](task-01-backend-scaffolding.md) | Go module, config, Chi router, middleware |
| 2 | [task-02-redis-cache.md](task-02-redis-cache.md) | Redis cache wrapper with JSON serialization |
| 3 | [task-03-nflverse-client.md](task-03-nflverse-client.md) | Parquet download, Apache Arrow parsing, caching |
| 4 | [task-04-aggregation-team-stats.md](task-04-aggregation-team-stats.md) | Team comparison stats (EPA, success rate, etc.) |
| 5 | [task-05-aggregation-play-stats.md](task-05-aggregation-play-stats.md) | Play-by-play processing, quarter/drive grouping |
| 6 | [task-06-http-handlers.md](task-06-http-handlers.md) | All REST endpoints + route registration |

## Frontend Tasks

| Task | File | Description |
|------|------|-------------|
| 7 | [task-07-frontend-api-client.md](task-07-frontend-api-client.md) | TypeScript types, ky client, React Query hooks |
| 8 | [task-08-season-week-pages.md](task-08-season-week-pages.md) | Season grid, Week list, GameCard component |
| 9 | [task-09-game-detail-stats.md](task-09-game-detail-stats.md) | GamePage layout, Stats tab, Recharts |
| 10 | [task-10-game-detail-pbp.md](task-10-game-detail-pbp.md) | Play-by-Play tab, quarter accordions |
| 11 | [task-11-game-detail-notes.md](task-11-game-detail-notes.md) | Notes tab: weather, injuries, vegas lines |

## Infrastructure Tasks

| Task | File | Description |
|------|------|-------------|
| 12 | [task-12-cicd-deployment.md](task-12-cicd-deployment.md) | GitHub Actions (API, Frontend, Lint), Fly.io |
| 13 | [task-13-docker-compose-local.md](task-13-docker-compose-local.md) | Local dev stack, Makefile, README |
| 14 | [task-14-polish-edge-cases.md](task-14-polish-edge-cases.md) | Error boundaries, skeletons, a11y, performance |

---

## Suggested Execution Order

### Phase 1: Backend Foundation (can parallelize)
1. Task 1 → Task 2 → Task 3 (sequential dependencies)
2. Task 4 & Task 5 (parallel, after Task 3)
3. Task 6 (after Task 4 & 5)

### Phase 2: Frontend Foundation (can parallelize with Phase 1 after Task 1)
4. Task 7 (after Task 1 for types)
5. Task 8 (after Task 7)
6. Task 9, 10, 11 (parallel, after Task 8)

### Phase 3: Integration & Deploy
7. Task 12 (CI/CD)
8. Task 13 (Local dev polish)
9. Task 14 (Final polish)

---

## Running a Single Task

```bash
# Read the task file
cat docs/superpowers/plans/tasks/task-01-backend-scaffolding.md

# Work on it, then commit
git add ...
git commit -m "feat: task 1 - backend scaffolding"
```

## Subagent-Driven Execution

For each task, you can use:
```bash
# In a new session, with the task file as context
opencode run --skill subagent-driven-development --task "Implement task-01-backend-scaffolding.md"
```

Or use `executing-plans` skill for native execution.
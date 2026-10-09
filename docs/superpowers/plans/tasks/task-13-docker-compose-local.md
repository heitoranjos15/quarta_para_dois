# Task 13: Docker Compose & Local Dev Polish

**Files:**
- Modify: `docker-compose.yml`
- Create: `backend/Dockerfile` (already exists)
- Create: `frontend/Dockerfile.dev` (already exists)
- Create: `Makefile` (optional convenience)

---

## Steps

- [ ] **Step 1: Verify `docker-compose up -d` starts API + Redis + Frontend**

```bash
docker-compose up -d
# Check logs
docker-compose logs -f api
docker-compose logs -f frontend
```

- [ ] **Step 2: Test full flow** — browser to `localhost:5173` → API → Redis → nflverse
  - Visit http://localhost:5173
  - Navigate Season → Week → Game
  - Verify Stats, PBP, Notes tabs load
  - Check Redis: `docker-compose exec redis redis-cli KEYS "*"`

- [ ] **Step 3: Add/Update README.md** with:
  - Architecture diagram
  - Quick start commands
  - Deployment guide
  - Environment variables reference

- [ ] **Step 4: Add Makefile** for common commands

```makefile
.PHONY: dev test lint build

dev:
	docker-compose up -d

dev-logs:
	docker-compose logs -f

backend-test:
	cd backend && go test ./...

backend-lint:
	cd backend && golangci-lint run

frontend-test:
	cd frontend && npm run test

frontend-lint:
	cd frontend && npm run lint

frontend-build:
	cd frontend && npm run build

test: backend-test frontend-test
lint: backend-lint frontend-lint
```

- [ ] **Step 5: Create feature branch**
```bash
git checkout -b feat/docker-compose-local
```

- [ ] **Step 6: Commit**
```bash
git add docker-compose.yml Makefile README.md
git commit -m "chore: docker compose polish, makefile, readme updates"
```
# Task 12: CI/CD & Deployment Config

**Files:**
- Create: `.github/workflows/api.yml`
- Create: `.github/workflows/frontend.yml`
- Create: `.github/workflows/lint.yml`
- Create: `fly.toml` (Fly.io config)

**Interfaces:**
- Produces: Automated deployments on push to main

---

## Steps

- [ ] **Step 1: API workflow** (`.github/workflows/api.yml`)

```yaml
name: Deploy API
on:
  push:
    branches: [main]
    paths: ['backend/**']
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --config backend/fly.toml --dockerfile backend/Dockerfile
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

- [ ] **Step 2: Frontend workflow** (`.github/workflows/frontend.yml`)

```yaml
name: Deploy Frontend
on:
  push:
    branches: [main]
    paths: ['frontend/**']
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20', cache: 'npm', cache-dependency-path: 'frontend/package-lock.json' }
      - run: cd frontend && npm ci && npm run build
      - uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./frontend/dist
          publish_branch: gh-pages
```

- [ ] **Step 3: Lint workflow** (`.github/workflows/lint.yml`)

```yaml
name: Lint
on: [push, pull_request]
jobs:
  backend-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: cd backend && golangci-lint run ./...
  frontend-lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20', cache: 'npm', cache-dependency-path: 'frontend/package-lock.json' }
      - run: cd frontend && npm ci && npm run lint
```

- [ ] **Step 4: Fly.io config** (`fly.toml`)

```toml
app = "quarta-para-dois-api"
primary_region = "iad"

[build]
  dockerfile = "backend/Dockerfile"

[env]
  PORT = "8080"
  ENV = "production"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
```

- [ ] **Step 5: Configure Fly.io app**
```bash
flyctl launch --name quarta-para-dois-api --dockerfile backend/Dockerfile
flyctl secrets set REDIS_ADDR=... REDIS_PASSWORD=... GITHUB_TOKEN=...
```

- [ ] **Step 6: Test deployments**

- [ ] **Step 7: Create feature branch**
```bash
git checkout -b feat/cicd-deployment
```

- [ ] **Step 8: Commit**
```bash
git add .github/workflows/ fly.toml
git commit -m "ci: add github actions for api, frontend, lint"
```
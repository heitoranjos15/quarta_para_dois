# Task 1: Backend Scaffolding & Config

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

---

## Steps

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
# Task 2: Redis Cache Wrapper

**Files:**
- Create: `backend/internal/nflverse/cache.go`
- Create: `backend/internal/nflverse/cache_test.go`

**Interfaces:**
- Consumes: `config.Config` (RedisAddr, RedisPassword)
- Produces: `Cache` interface with `Get(key string) (any, bool)`, `Set(key string, val any, ttl time.Duration)`, `Delete(key string)`

---

## Steps

- [ ] **Step 1: Write failing test**

```go
// backend/internal/nflverse/cache_test.go
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

```go
// backend/internal/nflverse/cache.go
type Cache interface {
    Get(key string) (any, bool)
    Set(key string, val any, ttl time.Duration)
    Delete(key string)
}

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr, password string) *RedisCache { ... }
```

- [ ] **Step 4: Run test**
Run: `cd backend && go test ./internal/nflverse/... -v`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/internal/nflverse/cache.go backend/internal/nflverse/cache_test.go
git commit -m "feat: redis cache wrapper"
```
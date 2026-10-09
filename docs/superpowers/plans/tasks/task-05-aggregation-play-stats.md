# Task 5: Aggregation Layer (Play-by-Play Processing)

**Files:**
- Create: `backend/internal/aggregation/play_stats.go`
- Create: `backend/internal/aggregation/play_stats_test.go`

**Interfaces:**
- Consumes: `[]Play` (single game)
- Produces: `[]PlayDetail` for UI table, grouped by quarter

---

## Steps

- [ ] **Step 1: Write failing test**

```go
// backend/internal/aggregation/play_stats_test.go
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

Define output structs:
```go
type PlayDetail struct {
    PlayID         int
    Quarter        int
    Down           int
    Distance       int
    Yardline       int
    PlayType       string
 midst
    Description    string
    EPA            float64
    Passer         string
    Receiver       string
    Rusher         string
    Tackler        string
    IsScoringPlay  bool
    DriveID        int
    DrivePlayCount int
    DriveYards     int
    DriveResult    string
    DriveTOP       time.Duration
}

type QuarterPlays struct {
    Quarter int
    Plays   []PlayDetail
    Drives  []DriveSummary
}
```

Implementation:
- Filter PBP to `game_id`
- Group by quarter (`qtr`)
- For each play: extract description, EPA, players (passer, receiver, rusher, tackler), play type, down/distance, yardline
- Group by drive: compute drive summary (start, end, result, plays, yards, TOP)
- Sort plays by `play_id` within quarter

- [ ] **Step 4: Run test** → PASS

- [ ] **Step 5: Create feature branch**
```bash
git checkout -b feat/aggregation-play-stats
```

- [ ] **Step 6: Commit**
```bash
git add backend/internal/aggregation/play_stats.go backend/internal/aggregation/play_stats_test.go
git commit -m "feat: play-by-play processing with drive grouping"
```
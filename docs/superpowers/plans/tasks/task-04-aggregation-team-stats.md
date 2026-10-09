# Task 4: Aggregation Layer (Team Stats)

**Files:**
- Create: `backend/internal/aggregation/team_stats.go`
- Create: `backend/internal/aggregation/team_stats_test.go`

**Interfaces:**
- Consumes: `[]Play` (filtered to single game)
- Produces: `TeamComparisonStats` struct with offense/defense metrics

---

## Steps

- [ ] **Step 1: Write failing test**

```go
// backend/internal/aggregation/team_stats_test.go
func TestComputeTeamStats(t *testing.T) {
    plays := loadTestPlays("2024_01_LAC_DEN") // fixture
    stats := aggregation.ComputeTeamStats(plays)
    assert.Equal(t, "LAC", stats.HomeTeam)
    assert.Equal(t, "DEN", stats.AwayTeam)
    assert.Greater(t, stats.HomeOffense.EPA_per_play, 0.0)
    assert.Greater(t, stats.AwayDefense.SuccessRate, 0.0)
}
```

- [ ] **Step 2: Run test**
Run: `cd backend && go test ./internal/aggregation/... -v`
Expected: FAIL

- [ ] **Step 3: Implement team_stats.go**

Define output struct:
```go
type TeamComparisonStats struct {
    HomeTeam string
    AwayTeam string
    HomeOffense  TeamUnitStats
    HomeDefense  TeamUnitStats
    AwayOffense  TeamUnitStats
    AwayDefense  TeamUnitStats
}

type TeamUnitStats struct {
    Plays            int
    EPA_per_play     float64
    SuccessRate      float64
    PassRate         float64
    RushEPA          float64
    PassEPA          float64
    RedZonePct       float64
    ThirdDownPct     float64
    FourthDownPct    float64
    TimeOfPossession time.Duration
    ExplosivePlays   int
    Turnovers        int
}
```

Compute from PBP (filter by `posteam`):
- EPA per play (offense/defense, pass/rush)
- Success rate (epa > 0)
- Red zone efficiency (plays inside 20)
- 3rd/4th down conversion
- Time of possession (sum of `game_seconds_remaining` diff)
- Explosive plays (20+ pass, 10+ rush)
- Turnover worthy plays

- [ ] **Step 4: Run test**
Run: `cd backend && go test ./internal/aggregation/... -v`
Expected: PASS

- [ ] **Step 5: Commit**
```bash
git add backend/internal/aggregation/
git commit -m "feat: team comparison stats aggregation"
```
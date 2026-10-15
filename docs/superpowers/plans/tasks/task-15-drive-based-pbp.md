# Task 15: Drive-Based Play-by-Play with Quarter Expansion

**Files:**
- Create: `backend/internal/aggregation/drive_stats.go`
- Create: `backend/internal/aggregation/drive_stats_test.go`
- Modify: `backend/internal/handlers/games.go` (add `/drives` endpoint)
- Create: `frontend/src/api/types.ts` (add drive types)
- Create: `frontend/src/api/client.ts` (add getDrives function)
- Create: `frontend/src/hooks/useDrives.ts`
- Create: `frontend/src/components/pbp/DriveList.tsx`
- Create: `frontend/src/components/pbp/DriveItem.tsx`
- Modify: `frontend/src/pages/GamePage.tsx` (update PlayByPlay tab)

**Interfaces:**
- Backend: `GET /api/games/{game_id}/drives` → `DriveSummaryResponse`
- Frontend: `useDrives(gameId)` hook

---

## Backend: Drive Aggregation

### Drive End Types (end reasons)
- `PUNT` - punt play
- `DOWNS` - 4th down not converted
- `INTERCEPTION` - interception
- `FUMBLE_LOST` - fumble lost
- `FIELD_GOAL` - field goal attempt (made/missed)
- `EXTRA_POINT` - extra point attempt
- `TOUCHDOWN` - touchdown
- `SAFETY` - safety
- `END_OF_HALF` - end of quarter/half
- `END_OF_GAME` - end of game
- `OTHER` - other (kickoff return TD, etc.)

### Response Structure

```json
{
  "data": {
    "quarters": {
      "1": {
        "drives": [
          {
            "drive_id": 1,
            "team": "LAC",
            "start_yardline": 25,
            "end_yardline": 0,
            "plays": 8,
            "yards": 75,
            "time_of_possession": "3:45",
            "result": "TOUCHDOWN",
            "end_play_description": "Herbert pass to Williams for 15 yards TD",
            "end_play_type": "TOUCHDOWN",
            "plays": [...]
          }
        ],
        "total": {
          "home": { "drives": 3, "points": 7, "time_of_possession": "8:30", "yards": 120 },
          "away": { "drives": 3, "points": 0, "time_of_possession": "6:30", "yards": 85 }
        }
      }
    },
    "overtime": { ... }
  },
  "meta": { "cached": true, "season": 2024 }
}
```

### Implementation Steps

1. **Create `drive_stats.go`** with:
   - `DriveEndReason` enum/constants
   - `DriveDetail` struct with all plays
   - `QuarterDriveSummary` with drives + team totals
   - `ComputeDrives(plays []Play) map[int]QuarterDriveSummary` function
   - Logic to group plays by drive_id, determine end reason from last play

2. **Add handler in `games.go`**:
   - `GET /api/games/{game_id}/drives`
   - Reuse existing PBP data, compute drives

3. **Add tests** for drive computation

---

## Frontend: Drive List with Expandable Quarters

### Types to Add (`types.ts`)

```typescript
export interface DriveDetail {
  drive_id: number;
  team: string;
  start_yardline: number;
  end_yardline: number;
  plays_count: number;
  yards: number;
  time_of_possession: string;
  result: string;
  end_play_description: string;
  end_play_type: string;
  plays: PlayDetail[];
}

export interface QuarterDriveSummary {
  drives: DriveDetail[];
  total: {
    home: TeamDriveTotal;
    away: TeamDriveTotal;
  };
}

export interface TeamDriveTotal {
  drives: number;
  points: number;
  time_of_possession: string;
  yards: number;
}

export interface DriveSummaryResponse {
  quarters: Record<number, QuarterDriveSummary>;
  overtime?: QuarterDriveSummary;
}
```

### Components to Create

1. **`DriveList.tsx`** - Main container with quarter accordions
2. **`DriveItem.tsx`** - Individual drive with expandable plays
3. **`DriveHeader.tsx`** - Drive summary row (team, result, yards, TOP)

### Layout Design

```
┌─ 1st Quarter ────────────── 3 drives • LAC 7 - LV 0 ─────▼┐
│  ┌─ Drive 1: LAC • TD • 75 yds • 3:45 ──────────────────▼┤
│  │  • Kickoff: Davis returns 35 yds                        │
│  │  • Pass: Herbert to Dissly 13 yds                       │
│  │  • Run: Edwards 6 yds                                   │
│  │  • ...                                                  │
│  │  • Pass: Herbert to Williams 15 yds TD  ← END PLAY     │
│  └────────────────────────────────────────────────────────┘
│  ┌─ Drive 2: LV • PUNT • 9 yds • 2:10 ──────────────────▼┤
│  │  ...                                                   │
│  └────────────────────────────────────────────────────────┘
└────────────────────────────────────────────────────────────┘

┌─ 2nd Quarter ────────────── 4 drives • LAC 10 - LV 7 ────▼┐
│  ...
└────────────────────────────────────────────────────────────┘
```

### Drive Header Format
`{Team} • {Result} • {Yards} yds • {TOP}` with color-coded result badge

### Play List Format
Each play: `Q{quarter} {down}&{distance} {type} {description} (EPA: {epa})`

---

## Steps

- [ ] **Backend: Create drive_stats.go** with drive computation logic
- [ ] **Backend: Add unit tests** for drive aggregation
- [ ] **Backend: Add /drives endpoint** in games.go handler
- [ ] **Frontend: Add types** for DriveDetail, QuarterDriveSummary
- [ ] **Frontend: Add getDrives API client** function
- [ ] **Frontend: Add useDrives hook**
- [ ] **Frontend: Create DriveList component** with quarter accordions
- [ ] **Frontend: Create DriveItem component** with expandable plays
- [ ] **Frontend: Update GamePage PlayByPlay tab** to use DriveList
- [ ] **Test with real data** (2024_01_LV_LAC)
- [ ] **Frontend: Add Playwright test scenarios** for drive-based PBP
  - [ ] Test drive list loads and displays quarters
  - [ ] Test quarter accordions expand/collapse with team totals
  - [ ] Test drive items expand to show plays
  - [ ] Test drive result badges (TD, FG, PUNT, INT, etc.)
  - [ ] Test end play description displays
  - [ ] Test play-by-play tab switches to drive view
  - [ ] Test overtime quarter displays if present
- [ ] **Create feature branch**: `feat/drive-based-pbp`
- [ ] **Commit**: `feat: drive-based play-by-play with quarter expansion`
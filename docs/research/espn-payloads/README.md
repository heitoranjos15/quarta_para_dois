# ESPN / weather payload samples

Raw responses captured 2026-10-21 while scoping the Notes tab. Two reference games:

- **401873007 — DAL @ GB**, week 6, final (Lambeau Field, outdoor)
- **401873010 — NE @ CHI**, week 7, scheduled for Thu 2026-10-22 (Soldier Field, outdoor)

| File | Source | What to look at |
|------|--------|-----------------|
| `scoreboard-2026-week6.json` / `week7.json` | `site.api.espn.com/.../scoreboard?seasontype=2&week=N&dates=2026` | `events[].status.type.name` (`STATUS_FINAL` / `STATUS_SCHEDULED`), used to pick the landing week |
| `summary-401873007-DAL-at-GB-final.json` | `site.api.espn.com/.../summary?event=ID` | `gameInfo.officials` (Referee = Craig Wrolstad), `injuries`, `pickcenter` (open/close odds), `drives.previous[].plays[].text` ("was injured during the play") |
| `summary-401873010-NE-at-CHI-scheduled.json` | same | `gameInfo.weather` (forecast; only present before kickoff), `injuries`, `pickcenter`; no `officials` yet |
| `core-odds-*.json` | `sports.core.api.espn.com/.../competitions/ID/odds` | only DraftKings listed |
| `core-odds-movement-*.json` | `.../odds/{provider}/history/0/movement` | empty (`count: 0`) for every provider tried |
| `weather-openmeteo-lambeau-2026-10-18.json` | `api.open-meteo.com/v1/forecast` (hourly, past date) | hourly temp / precip / snow / `weather_code` / wind for the game window |
| `weather-openmeteo-archive-lambeau-2026-10-18.json` | `archive-api.open-meteo.com/v1/archive` | same data from the historical archive |

## nflverse (`nflverse/` folder)

| File | Source | What to look at |
|------|--------|-----------------|
| `games.csv` | `github.com/nflverse/nfldata/raw/master/data/games.csv` | `espn` column maps `2026_06_DAL_GB` → `401873007`; closing `spread_line` / `total_line` / moneylines; `referee`, `temp`, `wind` empty for 2026 so far |
| `injuries_2026.parquet` | `nflverse-data` release `injuries/` | official game-week report: `report_status` (Out / Doubtful / Questionable), `report_primary_injury`, `practice_status` |
| `play_by_play_2026.parquet` | `nflverse-data` release `pbp/` | `weather` string ("Clear Temp: 59° F, Humidity: 44%, Wind: SW 10 mph"), `desc` with "was injured during the play" / "has returned to the game", `field_goal_result` |

`schedules/schedules_2026.parquet` (the URL the backend uses today) returns **404 Not Found**.

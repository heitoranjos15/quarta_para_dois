# Quarta Para Dois API — integration guide

The full contract is in [`openapi.yaml`](./openapi.yaml) (OpenAPI 3.1). Open it in
[Swagger Editor](https://editor.swagger.io), import it into Postman, or generate a Swift client
with [swift-openapi-generator](https://github.com/apple/swift-openapi-generator).

- **Base URL:** `https://quarta-para-dois-api.fly.dev` (local: `http://localhost:8080`)
- **Auth:** none. All endpoints are `GET`.
- **Cold starts:** the server scales to zero, so the first request after a quiet period can take a few seconds. Use a generous timeout (~30s).

## Endpoints

| Endpoint | Returns |
|---|---|
| `GET /api/seasons` | `[Int]` seasons, ascending |
| `GET /api/seasons/{year}/weeks` | `[Int]` weeks (1–18 regular season, 19+ postseason) |
| `GET /api/seasons/{year}/weeks/{week}/games` | `[Game]` |
| `GET /api/games/{gameID}` | `Game` |
| `GET /api/games/{gameID}/notes` | `GameNotes` (venue) |
| `GET /api/games/{gameID}/stats` | `TeamComparisonStats` |
| `GET /api/games/{gameID}/plays?quarter=&drive=` | `[QuarterPlays]` |
| `GET /api/games/{gameID}/drives` | `DriveSummaryResponse` |
| `GET /api/teams/{abbr}/games?season=` | `[Game]` |
| `GET /health` | plain text `ok` |

`gameID` is the nflverse ID `YYYY_WW_AWAY_HOME`, e.g. `2026_06_WAS_SF`. Take it from `Game.game_id`; don't build it yourself.

## Response shape

```json
{ "data": <payload>, "meta": { "cached": false, "season": 2026 } }
```

Errors return a non-2xx status with `{ "error": "game not found" }`.

| Status | Meaning | Suggested UI |
|---|---|---|
| 400 | Malformed year / week / game ID | Bug on our side, log it |
| 404 on `/stats`, `/plays`, `/drives` | Play-by-play not published yet (scheduled/live games, and usually until the day after a game) | Empty state ("Stats available after the game") |
| 404 elsewhere | Week or game doesn't exist | Not found |
| 503 | Upstream (ESPN / nflverse) failed | Retry |

## Gotchas for Codable

1. **All keys are snake_case**, so `decoder.keyDecodingStrategy = .convertFromSnakeCase` works across every endpoint.
2. **Durations come in two formats.**
   - Fields ending in `_ns` (`time_of_possession_ns` on `/stats`, `drive_top_ns` and `top_ns` on `/plays`): `Int64` **nanoseconds** (`/ 1_000_000_000` gives seconds).
   - `/drives` `time_of_possession`: `"M:SS"` string (may be `""`).
3. **Rates are fractions, 0–1** (`success_rate`, `red_zone_pct`, `third_down_pct`...). Multiply by 100 for display.
4. **Nullable fields:** `home_score`, `away_score`, `epa`, and `QuarterPlays.drives` (null when no drive started that quarter).
5. **Optional fields:** `Game.roof` and `Game.surface` are left out of the JSON when unknown.
6. **Empty strings mean unknown:** `status` (`""`), `kickoff`, `espn_id`, `location`, `game_type`, player names on plays.
7. **`/drives` quarters are a dictionary keyed by string** (`"1"`…`"4"`, `"5"` = overtime). Decode as `[String: QuarterDriveSummary]` and sort the keys numerically.
8. **Enums should decode leniently.** Use a fallback `unknown` case for `status` and drive `result` so a new value doesn't fail the whole decode.
9. **`yardline`** is yards from the *opponent's* end zone: 75 = own 25, 20 = opponent 20.
10. **`kickoff`** is RFC 3339 UTC (`2026-10-12T00:20:00Z`). `.iso8601` works, but the field can be `""`, so decode it as a `String` first.

## Example: `Game`

```json
{
  "game_id": "2026_06_WAS_SF",
  "espn_id": "401772940",
  "season": 2026,
  "week": 6,
  "game_type": "REG",
  "kickoff": "2026-10-12T00:20:00Z",
  "status": "final",
  "home_team": "SF",
  "away_team": "WAS",
  "home_score": 27,
  "away_score": 17,
  "stadium": "Levi's Stadium",
  "location": "Santa Clara, CA",
  "indoor": false,
  "roof": "outdoors",
  "surface": "grass"
}
```

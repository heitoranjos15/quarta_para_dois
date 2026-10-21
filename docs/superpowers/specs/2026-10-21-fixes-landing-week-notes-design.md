# Design: correctness fixes, current-week landing, and the new Notes tab

Date: 2026-10-21 · Repo state: `e2da3d33` on `master`

## Intent

quarta_para_dois is a free NFL game analytics portal (Go API + React SPA). Its owner uses it to go back over the week's games: on Wednesday, the games just played. Today the most visible numbers are wrong (drive results, points, some game pages 404), the app opens on the wrong week, and the Notes tab is mostly placeholder data.

This work makes the existing numbers trustworthy, makes the app open on the week that matters, and turns Notes into a real pre-game and post-game briefing: who's out, who got hurt, who's refereeing, what the weather is doing, and how the betting line moved.

## Goals

1. Everything shown is correct, or clearly marked as unavailable.
2. Opening the app lands on the current week's games with no clicks.
3. Notes works for **every game, before kickoff and after**.
4. CI actually runs and blocks regressions.

## Non-goals

- Line-movement history between open and close: parked in [issue #7](https://github.com/heitoranjos15/quarta_para_dois/issues/7).
- Grass species (e.g. Kentucky bluegrass). No free source has it.
- New stats or pages beyond what's listed here.

## Anti-goals (failure even if it "works")

- **Showing placeholder data as real.** No hard-coded "Clear", no 0% humidity, no "No injuries reported" when we simply don't have the data. Missing → hide the field or say "Not available" / "Not announced yet".
- **Using any paid service.** Project rule: everything free (ESPN public APIs, nflverse, Open-Meteo).

## Constraints

- Free data only: ESPN public site/core APIs (no key), nflverse GitHub releases and `nflverse/nfldata`, Open-Meteo (no key).
- Stack stays as is: Go 1.26 + chi, React 18 + TanStack Query + Tailwind; API on Fly.io, frontend on GitHub Pages.
- Default branch is `master`.
- Sample payloads for every source below are in `docs/research/espn-payloads/` (see its README). The reference games are **DAL @ GB** (`2026_06_DAL_GB`, ESPN `401873007`, final) and **NE @ CHI** (`2026_07_NE_CHI`, ESPN `401873010`, scheduled Thu 2026-10-22).

## Build order

Three parts, each with its own build plan, built in this order. The order was proposed by me and not objected to.

1. Fix what's wrong now (includes CI).
2. Current-week landing.
3. Notes tab.

---

## Part 1: Fix what's wrong now

Full findings with file/line references are in the review (summarized here). Each bug fix starts with a failing test.

### 1.1 Game identity: one canonical ID, mapped to ESPN *(user asked for this to be in the plan)*

**Problem:** the week list builds game IDs from ESPN abbreviations (`WSH`, `LAR`), but stats/plays/drives look games up in nflverse (`WAS`, `LA`). Every Washington or Rams game page 404s. `TeamLogo.tsx` also lacks `LA`.

**Fix:** nflverse's `games.csv` (`https://github.com/nflverse/nfldata/raw/master/data/games.csv`) has an **`espn` column mapping each game to its ESPN event ID** (`2026_06_DAL_GB` ↔ `401873007`). Use it:

- The canonical game ID everywhere (URLs, API, frontend) is the **nflverse `game_id`**.
- When ESPN data is needed, look up the ESPN event ID through `games.csv` and call ESPN by event ID, never by rebuilding IDs from abbreviations.
- ESPN scoreboard events are joined to nflverse games by ESPN event ID, never by team abbreviation.
- Frontend team colors/logos use nflverse abbreviations (`WAS`, `LA`, …).

### 1.2 Schedule source

The schedules parquet the backend uses (`.../releases/download/schedules/schedules_%d.parquet`) returns **404**. Remove it.

*User decision:* **ESPN is the source for games.** The week's game list, teams, kickoff time, live/final scores, game status, venue and indoor flag all come from the ESPN scoreboard (`seasontype=2` regular season, `3` playoffs).

`games.csv` is used **only as a lookup table**: ESPN event ID → nflverse `game_id` (for stats/plays/drives/injuries), plus `roof`, `surface` and `stadium_id`. Its scores and schedule are not used.

*User decision:* **no local files.** The app fetches `games.csv` from GitHub at runtime and keeps it only in cache (memory/Redis, refreshed periodically, e.g. daily). It is never committed, bundled into the build or read from disk. The copy in `docs/research/espn-payloads/nflverse/` is a research sample only. If an ESPN event has no `games.csv` match yet (e.g. a game added very recently), the game still lists from ESPN data and its nflverse-backed tabs say "Not available yet".

### 1.3 Drive results and points

- Check punt / field goal / extra point **before** the "turnover on downs" rule (punts and FGs happen on 4th down, so they're currently all "DOWNS" and 0 points).
- A field goal counts 3 points only when `field_goal_result == "made"`; otherwise it's a missed FG.
- Defensive or return touchdowns (pick-six, fumble return) end the drive as a turnover and credit no points to the offense.
- Group drives across quarter breaks: a drive is one drive even when it spans Q1→Q2. Per-quarter views may still show it, but totals must count it once.
- Same rules in `play_stats.go` (or delete the duplicated logic and reuse one implementation).
- Fix the test fixture that hid this (punt on `Down: 1`) and add cases for 4th-down punt, made/missed FG, pick-six and a cross-quarter drive.

### 1.4 Team stats

- Only count real scrimmage plays (pass/run) in per-play rates (EPA/play, success rate, pass rate). Today timeouts, quarter-end rows, kickoffs, PATs and no-plays dilute them.
- Rename or recompute "Red Zone %": the label should mean touchdowns per red-zone trip.
- Fill both sides of each comparison row on the game page (`StatSection` passes `away={0}` everywhere and `home={0}` for time of possession).

### 1.5 Other data bugs

- Team schedule endpoint: IDs come out as `YYYY_00_X_Y` and the season filter is a no-op. Fix via 1.1.
- Playoffs: support postseason weeks (ESPN `seasontype=3`). Week lists must be right per season (17 weeks for 2002–2020, 18 from 2021; nflverse `week` + `game_type` is the source of truth).
- Live scores: never serve a cached score for a game in progress. A week's ESPN scoreboard is not cached while any of its games is live (see 3.7).
- Parquet nulls must be read as null, not 0 (e.g. missing EPA).
- An invalid game ID returns 400/404, not a silent fallback to the current season.
- Winner highlighting: ties and unplayed games highlight nobody.
- Team logo initials must be unique (NYG/NYJ, LAC/LA/LV).

### 1.6 Performance and production

- Don't decode a whole season of play-by-play from Redis per request. Keep parsed data in memory indexed by game ID and/or cache per game. Read only the columns needed, and look up column indexes once per batch.
- Production needs a reachable cache, or must work well without one. `REDIS_ADDR` isn't set on Fly today and cache errors are swallowed silently. Must remain free. The project docs already point at Upstash, so production Redis is Upstash's free tier via Fly secrets (`REDIS_ADDR`, `REDIS_PASSWORD`). *My call:* cache per game rather than per season so values stay small (free-tier record and bandwidth limits), and keep the API working on in-memory data when Redis is down. The Fly machine is small, so in-memory data must be bounded (e.g. only the seasons/games recently requested).
- Server: timeouts on the HTTP server, return/log the `ListenAndServe` error, remove `DEBUG` printfs, use `slog` (JSON when `ENV=production`), inject upstream clients into handlers via constructor closures (internal wiring only, unrelated to the endpoint split in 3.8; drop the context-injection middleware), delete unused `middleware/logging.go` (it pulls in chi v4).
- One Redis client shared through the `Cache` interface (the ESPN client opens its own today).

### 1.7 CI and deploy

- Workflows trigger on `master`, not `main`.
- `api.yml`: correct the `fly.toml` path (it's at the repo root).
- Go version in CI matches `go.mod` (1.26); install golangci-lint via its action; add a `go test ./...` job; `gofmt` the codebase.
- Frontend: add `eslint.config.js` (ESLint 9 flat config), add `@types/node` so `tsc -b` / `npm run build` pass, keep vitest away from the Playwright `tests/` folder, add a CI test step.
- GitHub Pages: set Vite `base` to the repo path (unless a custom domain is configured) and make deep links survive a refresh (`404.html` fallback or hash routing; builder's choice).
- Dockerfile: no `go mod tidy` during the build.

### 1.8 Cleanup (when touching the area)

Duplicate `DriveSummaryResponse`; unused `Overtime` field (overtime is `quarters[5]`); `Play.YardsToGoal` duplicate tag; unreachable `case ESPNScore`. Adding `json:"snake_case"` tags to Go response structs lets most of the frontend mapping code in `types.ts` go away. Allowed, not required.

---

## Part 2: Current-week landing *(user's rule)*

**Rule:** opening `/` goes straight to `/season/{season}/week/{week}`, where the week is **the latest week whose first game has finished**. That week stays the landing week until the first game of the next week finishes.

Example: on Wed 2026-10-21 it's week 6. After Thursday's NE @ CHI (week 7's first game) goes final, it's week 7.

- The backend decides, from real game status on the ESPN scoreboard (`status.type.name == STATUS_FINAL`). The browser's calendar math (`getCurrentNFLWeek`, which says week 8 today) is removed. Suggested shape: `GET /api/current-week` → `{ season, week, season_type }`.
- *My call:* the same rule continues through the playoffs (Wild Card → Super Bowl). Playoff weeks use nflverse numbering in URLs (weeks 19–22 for 18-week seasons) and are labeled Wild Card / Divisional / Conference / Super Bowl in the UI. Live game status comes from ESPN's scoreboard (`seasontype=3` for playoffs), joined by ESPN event ID.
- *My call:* in the offseason (after the Super Bowl until the next opener's first game finishes), land on the last week played.
- `/season/{season}` without a week also lands on that season's current week (or week 1 / the final week for past seasons; builder's choice).
- The week strip on the week page stays; manual navigation is unchanged.

---

## Part 3: Notes tab

Works for **every game, before kickoff and after**. Each section hides or says "Not available" when its data is missing (see anti-goals).

### 3.1 Injury report *(user: use nflverse)*

- Source: **nflverse injuries** `https://github.com/nflverse/nflverse-data/releases/download/injuries/injuries_{season}.parquet`, rows for the game's `season`, `week` and both teams.
- *User decision:* show **every player on the injury report** for the game week, grouped by team, with position, injury (`report_primary_injury`, else `practice_primary_injury`) and status. Game status is `report_status` (Out / Doubtful / Questionable). Players listed only for practice show their `practice_status` (e.g. "Limited participation"). Sort by severity (Out first).
- The section's title should not be "Ruled out", since most listed players aren't. Use "Injury report".
- Same source before and after kickoff. The report fills in during the week (week 7 already has rows on Wednesday) and is preserved per week afterward.
- Before the report exists: "Injury report not published yet".
- nflverse injury data starts in 2009; for older seasons the section says "Not available".
- nflverse injury rows use nflverse team abbreviations, which match the canonical game ID (1.1).

### 3.2 Injured during the game

- Source: play-by-play `desc` (nflverse pbp; ESPN play text has the same wording).
- Patterns: `TEAM-NN-F.Last was injured during the play.` and `** Injury Update: TEAM-NN-F.Last has returned to the game.`
- Show per team: player, quarter + game clock of the injury, and whether/when they returned.
- Before kickoff: section hidden. If pbp for a just-finished game isn't in nflverse yet, use ESPN's play text.

### 3.3 Referee

- Source: ESPN game summary `gameInfo.officials`, position `Referee` (full crew optional).
- Before announcement (ESPN returns no officials): "Not announced yet".

### 3.4 Venue: indoors and surface *(user asked)*

- Stadium name.
- Roof: nflverse `roof`, shown as Outdoors / Indoors (`dome`) / Retractable roof, closed / Retractable roof, open. Fallback when blank: ESPN `venue.indoor`.
- Surface: nflverse `surface` (grass, fieldturf, sportturf, matrixturf, a_turf, astroturf), shown in readable form ("Natural grass", "FieldTurf", …). Fallback when blank (~1 in 10 games, e.g. DAL @ GB): ESPN `gameInfo.venue.grass` → "Natural grass" / "Artificial turf".

### 3.5 Weather *(user: Open-Meteo, including a forecast before kickoff)*

- Outdoor games and open retractable roofs only; indoor games show "Indoors" and no weather.
- *My call:* retractable-roof games whose roof status isn't known yet (before kickoff) show the forecast with a note "Roof decision pending".
- Source: **Open-Meteo**, hourly, at the stadium's coordinates, for the game window (kickoff to roughly 3.5 hours after kickoff, or to the actual end time when known).
  - Before kickoff: forecast API (`api.open-meteo.com/v1/forecast`).
  - After the game: the same API's past days, or the archive API (`archive-api.open-meteo.com/v1/archive`) for older games.
- Each hour: condition (sunny/clear, cloudy, rain, storm, snow, fog) from `weather_code`, temperature, precipitation / snowfall, wind and gusts.
- *User decision:* US units throughout: °F, mph, inches.
- Summary line: kickoff conditions plus any change during the game ("Clear 56°F → 49°F, wind gusts to 17 mph").
- *My call:* stadium coordinates come from a static table keyed by nflverse `stadium_id`, covering current stadiums and international venues. If an ID is missing, fall back to Open-Meteo's free geocoding by city.

### 3.6 Odds

- Source: ESPN core odds `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/events/{espnId}/competitions/{espnId}/odds` (DraftKings is the only provider listed).
- Show **open → close** for spread, moneyline and total. Before kickoff, **open → current**.
- After the game: whether the favorite won (`moneylineWinner`) and covered (`spreadWinner`), and whether the total went over or under (derived from the final score).
- Line movement in between: out of scope (issue #7).

### 3.7 Freshness

*User decision:* keep TTLs high. Cache upstream data by game state:

- **Finished games:** long TTL (*my call:* 30 days). Their play-by-play, results, officials, closing odds and observed weather don't change.
- **Upcoming games:** shorter TTL so the injury report, odds and forecast stay current (*my call:* about 1 hour).
- **Live games:** *user decision:* **no caching for now.** While a game is in progress, its score, status and play data are fetched fresh on every request.
- **`games.csv` lookup:** about a day.

### 3.8 API shape

*User decision:* endpoints are organized by **what the data is, not which upstream it comes from**. One endpoint can merge ESPN + nflverse + Open-Meteo data. Injuries and weather get their own endpoints.

- `GET /api/games/{gameId}`: the game, merged from all sources. Teams, kickoff, score, status (ESPN); nflverse/ESPN IDs; venue (stadium, roof/indoors, surface); head referee (and crew); odds open → close/current with winner/cover/over-under.
- `GET /api/games/{gameId}/injuries`: the nflverse injury report (3.1) and in-game injuries (3.2).
- `GET /api/games/{gameId}/weather`: the Open-Meteo hourly forecast/observations for the game window (3.5), or an "indoors" marker.
- Existing `/stats`, `/plays`, `/drives` stay. `/notes` is replaced by the endpoints above.
- `GET /api/current-week` (Part 2).
- Each endpoint caches by its own data's freshness (3.7), so a slow source (e.g. weather) doesn't block the game header.

---

## Left to the builder

- Package/file layout for new upstream clients (nflverse injuries, `games.csv`, Open-Meteo, ESPN summary/odds); the endpoint split is fixed in 3.8.
- In-memory vs Redis split and fine-tuning the TTLs in 3.7, within "free" and "correct".
- The Notes tab's visual layout (sections in the order above unless there's a reason not to).
- How deep links survive refresh on GitHub Pages.
- Retry/backoff details for upstream calls.

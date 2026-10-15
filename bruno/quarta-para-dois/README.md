# Quarta Para Dois - Bruno API Collection

This folder contains a Bruno API collection for testing the NFL Game Analytics API.

## Setup

1. Install [Bruno](https://www.usebruno.com/)
2. Open Bruno and select "Open Collection"
3. Navigate to `bruno/quarta-para-dois` (the folder containing `bruno.json`)

## Endpoints

### Health Check
- `health.bru` - `GET /health` - Server health check

### Seasons
- `seasons-list.bru` - `GET /api/seasons` - List all available seasons (1999-present)
- `seasons-weeks.bru` - `GET /api/seasons/{year}/weeks` - List weeks in a season
- `seasons-games.bru` - `GET /api/seasons/{year}/weeks/{week}/games` - List games in a specific week

### Games
- `games-detail.bru` - `GET /api/games/{gameID}` - Get game summary
- `games-stats.bru` - `GET /api/games/{gameID}/stats` - Get team comparison stats (offense/defense EPA, success rate, etc.)
- `games-plays.bru` - `GET /api/games/{gameID}/plays` - Get play-by-play data
- `games-plays-quarter.bru` - `GET /api/games/{gameID}/plays?quarter={1-4}` - Get plays filtered by quarter
- `games-plays-drive.bru` - `GET /api/games/{gameID}/plays?drive={N}` - Get plays filtered by drive
- `games-notes.bru` - `GET /api/games/{gameID}/notes` - Get game notes (weather, injuries, vegas lines)

### Teams
- `teams-games.bru` - `GET /api/teams/{abbr}/games?season={year}` - Get all games for a team in a season

## Variables

Replace the following variables in requests (Bruno supports variables):
- `{{gameID}}` - Game ID format: `2024_01_LAC_DEN` (season_week_home_away)
- `{{year}}` - Season year (e.g., `2024`)
- `{{week}}` - Week number (1-18)
- `{{abbr}}` - Team abbreviation (e.g., `LAC`, `DEN`, `KC`)
- `{{quarter}}` - Quarter number (1-4)
- `{{drive}}` - Drive number

## Example Game IDs
- `2024_01_LAC_DEN` - Week 1, Chargers vs Broncos
- `2024_02_KC_BAL` - Week 2, Chiefs vs Ravens
- `2023_18_SF_LAR` - Week 18, 49ers vs Rams

## Notes

- First request for a season downloads ~180MB Parquet from nflverse-data (takes 30-60s)
- Subsequent requests for the same season are cached in Redis (<200ms)
- Requires Redis running on localhost:6379 (or set REDIS_ADDR env var)
- Optional: Set GITHUB_TOKEN for higher GitHub API rate limits
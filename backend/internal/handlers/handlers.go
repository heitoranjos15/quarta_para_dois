package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"quarta-para-dois/backend/internal/aggregation"
	"quarta-para-dois/backend/internal/nflverse"
)

const (
	nflverseClientKey = "nflverseClient"
	espnClientKey     = "espnClient"
)

type espnClientInterface interface {
	GetSchedule(ctx context.Context, season, week int) ([]nflverse.Game, error)
	GetGame(ctx context.Context, season, week int, gameID string) (*nflverse.Game, error)
	GetTeamSchedule(ctx context.Context, teamAbbr string, season int) ([]nflverse.Game, error)
}

func getESPNClient(r *http.Request) (espnClientInterface, bool) {
	client, ok := r.Context().Value(espnClientKey).(espnClientInterface)
	fmt.Printf("DEBUG getESPNClient: key=%q, ok=%v, client=%v\n", espnClientKey, ok, client)
	return client, ok
}

type APIResponse[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Cached bool `json:"cached"`
	Season int  `json:"season"`
}

type GameNotes struct {
	Weather  *WeatherInfo `json:"weather,omitempty"`
	Injuries []InjuryInfo `json:"injuries,omitempty"`
	Vegas    *VegasInfo   `json:"vegas,omitempty"`
}

type WeatherInfo struct {
	Temperature int    `json:"temperature"`
	Wind        int    `json:"wind"`
	Condition   string `json:"condition"`
	Stadium     string `json:"stadium"`
	Roof        string `json:"roof"`
}

type InjuryInfo struct {
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Detail   string `json:"detail"`
}

type VegasInfo struct {
	Spread float64 `json:"spread"`
	Total  float64 `json:"total"`
	HomeML int     `json:"home_ml"`
	AwayML int     `json:"away_ml"`
}

type nflverseClientInterface interface {
	GetPBP(season int) ([]nflverse.Play, error)
	GetSchedules(season int) ([]nflverse.Game, error)
	GetTeams() ([]nflverse.Team, error)
	GetRosters(season int) ([]nflverse.Roster, error)
}

func respond[T any](w http.ResponseWriter, data T, cached bool, season int, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := APIResponse[T]{
		Data: data,
		Meta: Meta{Cached: cached, Season: season},
	}
	json.NewEncoder(w).Encode(resp)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func extractSeason(gameID string) int {
	// Format: 2024_01_LAC_DEN
	parts := strings.Split(gameID, "_")
	if len(parts) > 0 {
		if season, err := strconv.Atoi(parts[0]); err == nil {
			return season
		}
	}
	return time.Now().Year()
}

func filterPlaysByGame(plays []nflverse.Play, gameID string) []nflverse.Play {
	var result []nflverse.Play
	for _, p := range plays {
		if p.GameID == gameID {
			result = append(result, p)
		}
	}
	return result
}

func SeasonsListHandler(w http.ResponseWriter, r *http.Request) {
	currentYear := time.Now().Year()
	if time.Now().Month() < 9 {
		currentYear--
	}

	seasons := make([]int, 0, currentYear-1998)
	for y := 1999; y <= currentYear; y++ {
		seasons = append(seasons, y)
	}

	respond(w, seasons, false, currentYear, http.StatusOK)
}

func SeasonsWeeksHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1999 {
		respondError(w, "invalid season", http.StatusBadRequest)
		return
	}

	// Regular season weeks (1-18) + playoffs
	maxWeek := 18
	if year >= 2021 {
		maxWeek = 18
	} else if year >= 1978 {
		maxWeek = 16
	}

	weeks := make([]int, maxWeek)
	for i := 1; i <= maxWeek; i++ {
		weeks[i-1] = i
	}

	respond(w, weeks, false, year, http.StatusOK)
}

func SeasonsGamesHandler(w http.ResponseWriter, r *http.Request) {
	yearStr := chi.URLParam(r, "year")
	weekStr := chi.URLParam(r, "week")

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1999 {
		respondError(w, "invalid season", http.StatusBadRequest)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 {
		respondError(w, "invalid week", http.StatusBadRequest)
		return
	}

	espnClient, ok := getESPNClient(r)
	if !ok {
		respondError(w, "internal error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG: Fetching schedule for %d week %d\n", year, week)
	ctx := r.Context()
	games, err := espnClient.GetSchedule(ctx, year, week)
	if err != nil {
		fmt.Printf("DEBUG: Error fetching schedule: %v\n", err)
		respondError(w, "failed to fetch schedules", http.StatusServiceUnavailable)
		return
	}

	fmt.Printf("DEBUG: Got %d games\n", len(games))
	respond(w, games, false, year, http.StatusOK)
}

func GameDetailHandler(nflClient nflverseClientInterface, espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season := extractSeason(gameID)

		ctx := r.Context()
		game, err := espnClient.GetGame(ctx, season, 0, gameID) // week=0 means search all weeks
		if err != nil {
			respondError(w, "failed to fetch game", http.StatusServiceUnavailable)
			return
		}

		if game == nil {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		respond(w, game, false, season, http.StatusOK)
	}
}

func GameStatsHandler(client nflverseClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season := extractSeason(gameID)

		plays, err := client.GetPBP(season)
		if err != nil {
			respondError(w, "failed to fetch play data", http.StatusServiceUnavailable)
			return
		}

		gamePlays := filterPlaysByGame(plays, gameID)
		if len(gamePlays) == 0 {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		stats := aggregation.ComputeTeamStats(gamePlays)
		respond(w, stats, true, season, http.StatusOK)
	}
}

func GamePlaysHandler(client nflverseClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season := extractSeason(gameID)

		// Optional query params for filtering
		quarterStr := r.URL.Query().Get("quarter")
		driveStr := r.URL.Query().Get("drive")

		plays, err := client.GetPBP(season)
		if err != nil {
			respondError(w, "failed to fetch play data", http.StatusServiceUnavailable)
			return
		}

		gamePlays := filterPlaysByGame(plays, gameID)
		if len(gamePlays) == 0 {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		// Filter by quarter if specified
		if quarterStr != "" {
			quarter, err := strconv.Atoi(quarterStr)
			if err == nil {
				var filtered []nflverse.Play
				for _, p := range gamePlays {
					if p.Quarter == quarter {
						filtered = append(filtered, p)
					}
				}
				gamePlays = filtered
			}
		}

		// Filter by drive if specified
		if driveStr != "" {
			drive, err := strconv.Atoi(driveStr)
			if err == nil {
				var filtered []nflverse.Play
				for _, p := range gamePlays {
					if p.Drive == drive {
						filtered = append(filtered, p)
					}
				}
				gamePlays = filtered
			}
		}

		details := aggregation.ProcessPlays(gamePlays)
		respond(w, details, true, season, http.StatusOK)
	}
}

func GameNotesHandler(nflClient nflverseClientInterface, espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := chi.URLParam(r, "gameID")
		season := extractSeason(gameID)

		ctx := r.Context()
		game, err := espnClient.GetGame(ctx, season, 0, gameID)
		if err != nil {
			respondError(w, "failed to fetch game", http.StatusServiceUnavailable)
			return
		}

		if game == nil {
			respondError(w, "game not found", http.StatusNotFound)
			return
		}

		notes := GameNotes{
			Weather: &WeatherInfo{
				Temperature: game.Temp,
				Wind:        game.Wind,
				Condition:   "Clear",
				Stadium:     game.Stadium,
				Roof:        game.Roof,
			},
			Injuries: []InjuryInfo{},
			Vegas: &VegasInfo{
				Spread: 0,
				Total:  0,
			},
		}

		respond(w, notes, false, season, http.StatusOK)
	}
}

func TeamGamesHandler(nflClient nflverseClientInterface, espnClient espnClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamAbbr := chi.URLParam(r, "abbr")
		seasonStr := r.URL.Query().Get("season")
		
		season := time.Now().Year()
		if time.Now().Month() < 9 {
			season--
		}
		if seasonStr != "" {
			if s, err := strconv.Atoi(seasonStr); err == nil {
				season = s
			}
		}
		
		ctx := r.Context()
		games, err := espnClient.GetTeamSchedule(ctx, teamAbbr, season)
		if err != nil {
			respondError(w, "failed to fetch team schedule", http.StatusServiceUnavailable)
			return
		}
		
		respond(w, games, false, season, http.StatusOK)
	}
}


package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"quarta-para-dois/backend/internal/config"
	"quarta-para-dois/backend/internal/handlers"
	"quarta-para-dois/backend/internal/middleware"
	"quarta-para-dois/backend/internal/nflverse"
)

const (
	nflverseClientKey = "nflverseClient"
	espnClientKey     = "espnClient"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	cache := nflverse.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword)
	nflClient := nflverse.NewClient(cache, cfg.GitHubToken)
	espnClient := nflverse.NewESPNClient(cfg.RedisAddr, cfg.RedisPassword)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Cors)

	// Add clients to request context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("DEBUG middleware: adding clients to context\n")
			ctx := r.Context()
			ctx = context.WithValue(ctx, nflverseClientKey, nflClient)
			ctx = context.WithValue(ctx, espnClientKey, espnClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		// Seasons
		r.Get("/seasons", handlers.SeasonsListHandler)
		r.Get("/seasons/{year}/weeks", handlers.SeasonsWeeksHandler)
		r.Get("/seasons/{year}/weeks/{week}/games", handlers.SeasonsGamesHandler)

		// Games
		r.Get("/games/{gameID}", handlers.GameDetailHandler(nflClient, espnClient))
		r.Get("/games/{gameID}/stats", handlers.GameStatsHandler(nflClient))
		r.Get("/games/{gameID}/plays", handlers.GamePlaysHandler(nflClient))
		r.Get("/games/{gameID}/notes", handlers.GameNotesHandler(nflClient, espnClient))

		// Teams
		r.Get("/teams/{abbr}/games", handlers.TeamGamesHandler(nflClient, espnClient))
	})

	log.Printf("Server starting on :%s", cfg.Port)
	_ = http.ListenAndServe(":"+cfg.Port, r)
}
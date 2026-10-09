package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"quarta-para-dois/backend/internal/config"
	"quarta-para-dois/backend/internal/handlers"
	"quarta-para-dois/backend/internal/middleware"
	"quarta-para-dois/backend/internal/nflverse"
)

type contextKey string

const nflverseClientKey contextKey = "nflverseClient"

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	cache := nflverse.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword)
	client := nflverse.NewClient(cache, cfg.GitHubToken)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Cors)

	// Add nflverse client to request context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, nflverseClientKey, client)
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
		r.Get("/games/{gameID}", handlers.GameDetailHandler(client))
		r.Get("/games/{gameID}/stats", handlers.GameStatsHandler(client))
		r.Get("/games/{gameID}/plays", handlers.GamePlaysHandler(client))
		r.Get("/games/{gameID}/notes", handlers.GameNotesHandler(client))

		// Teams
		r.Get("/teams/{abbr}/games", handlers.TeamGamesHandler(client))
	})

	log.Printf("Server starting on :%s", cfg.Port)
	_ = http.ListenAndServe(":"+cfg.Port, r)
}
package main

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"quarta-para-dois/backend/internal/config"
	"quarta-para-dois/backend/internal/handlers"
	"quarta-para-dois/backend/internal/middleware"
	"quarta-para-dois/backend/internal/nflverse"
)

func newLogger(env string, w io.Writer) *slog.Logger {
	if env == "production" {
		return slog.New(slog.NewJSONHandler(w, nil))
	}
	return slog.New(slog.NewTextHandler(w, nil))
}

func newServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "err", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.Env, os.Stdout)
	slog.SetDefault(logger)

	cache := nflverse.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword)
	nflClient := nflverse.NewClient(cache, cfg.GitHubToken)
	espnClient := nflverse.NewESPNClient(cfg.RedisAddr, cfg.RedisPassword)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestLogger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Cors)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		// Seasons
		r.Get("/seasons", handlers.SeasonsListHandler)
		r.Get("/seasons/{year}/weeks", handlers.SeasonsWeeksHandler)
		r.Get("/seasons/{year}/weeks/{week}/games", handlers.SeasonsGamesHandler(espnClient))

		// Games - specific routes first (chi matches in order)
		r.Get("/games/{gameID}/stats", handlers.GameStatsHandler(nflClient))
		r.Get("/games/{gameID}/plays", handlers.GamePlaysHandler(nflClient))
		r.Get("/games/{gameID}/drives", handlers.GameDrivesHandler(nflClient))
		r.Get("/games/{gameID}/notes", handlers.GameNotesHandler(nflClient, espnClient))
		r.Get("/games/{gameID}", handlers.GameDetailHandler(nflClient, espnClient))

		// Teams
		r.Get("/teams/{abbr}/games", handlers.TeamGamesHandler(nflClient, espnClient))
	})

	logger.Info("server starting", "port", cfg.Port)
	if err := newServer(":"+cfg.Port, r).ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

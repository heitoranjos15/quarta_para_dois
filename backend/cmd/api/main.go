package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v4"
	chimiddleware "github.com/go-chi/chi/v4/middleware"

	"quarta-para-dois/backend/internal/config"
	custommiddleware "quarta-para-dois/backend/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(custommiddleware.Cors)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("Server starting on :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
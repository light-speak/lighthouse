package net

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/net/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	setMiddlewares(r)
	setRoutes(r)
	return r
}

func setMiddlewares(r *chi.Mux) {
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.NoCache)
	r.Use(chiMiddleware.Heartbeat("/heartbeat"))
	r.Use(chiMiddleware.RequestLogger(&middleware.LogMiddleware{}))
	r.Use(chiMiddleware.Compress(5))
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(chiMiddleware.Throttle(env.LighthouseConfig.Server.Throttle))
	r.Use(middleware.AuthMiddleware)
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

func setRoutes(r *chi.Mux) {
	r.Options("/query", optionsHandler)
	r.Post("/query", graphQLHandler)
	r.Get("/query", graphQLHandler)
}

package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/light-speak/lighthouse/routers/health"
	"github.com/rs/cors"
)

const (
	ContentTypeJSON = "application/json"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   Config.CORSAllowOrigins,
		AllowedMethods:   Config.CORSAllowMethods,
		AllowedHeaders:   Config.CORSAllowHeaders,
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}).Handler)

	setMiddlewares(r)
	registerSystemRoutes(r)
	return r
}

func setMiddlewares(r *chi.Mux) {
	r.Use(middleware.Recoverer) // Recover from panics
	r.Use(middleware.RequestID) // Request ID
	r.Use(middleware.RealIP)    // Real IP
}

func registerSystemRoutes(r *chi.Mux) {
	r.Get(Config.HeartbeatPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.Get(Config.ReadinessPath, health.ReadinessHandler)
}

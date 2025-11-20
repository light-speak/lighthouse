package routers

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/light-speak/lighthouse/routers/log"
	"github.com/light-speak/lighthouse/routers/request"
	"github.com/rs/cors"
)

const (
	ContentTypeJSON = "application/json"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}).Handler)
	setMiddlewares(router)
	return router
}

func setMiddlewares(r *chi.Mux) {

	r.Use(middleware.Recoverer)                             // Recover from panics
	r.Use(middleware.RequestID)                             // Request ID
	r.Use(middleware.RealIP)                                // Real IP
	r.Use(request.RequestMiddleware)                        // Request middleware
	r.Use(middleware.NoCache)                               // No cache
	r.Use(middleware.Heartbeat(Config.HeartbeatPath))       // Heartbeat
	r.Use(middleware.RequestLogger(&log.LogMiddleware{}))   // Request logger
	r.Use(middleware.Compress(Config.CompressLevel))        // Compress
	r.Use(httprate.LimitByIP(Config.Throttle, time.Minute)) // Limit by IP
}

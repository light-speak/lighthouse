package middleware

import (
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(ContextMiddleware())
	return router
}

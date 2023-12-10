package middleware

import (
	"github.com/go-chi/chi"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(ContextMiddleware())
	return router
}

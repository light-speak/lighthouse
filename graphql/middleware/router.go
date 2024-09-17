package middleware

import (
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(ContextMiddleware())
	
	// 可以在这里添加更多的全局中间件
	// 例如: router.Use(cors.Handler(cors.Options{...}))
	
	return router
}

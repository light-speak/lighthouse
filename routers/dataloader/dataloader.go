package dataloader

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

type Dataloader interface {
	Name() string
	IsLoader() bool
	NewLoader(db *gorm.DB) Dataloader
}

var loaders = make(map[string]Dataloader)

func RegisterLoader(loader Dataloader) {
	loaders[loader.Name()] = loader
}

func GetLoader(name string) Dataloader {
	return loaders[name]
}

type LoaderContextKey struct {
	name string
}

func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			for _, loader := range loaders {
				l := loader.NewLoader(db)
				ctx = context.WithValue(ctx, LoaderContextKey{loader.Name()}, l)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLoaderFromCtx(ctx context.Context, name string) Dataloader {
	return ctx.Value(LoaderContextKey{name}).(Dataloader)
}

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

// Config holds dataloader configuration
type Config struct {
	// QueryTimeout is the timeout for database queries in dataloader
	// Default: 15 seconds
	QueryTimeout time.Duration
}

var config = &Config{
	QueryTimeout: 15 * time.Second,
}

func init() {
	// Allow override via environment variable (in seconds)
	timeout := utils.GetEnvInt("DATALOADER_QUERY_TIMEOUT", 15)
	config.QueryTimeout = time.Duration(timeout) * time.Second
}

// GetQueryTimeout returns the configured query timeout
func GetQueryTimeout() time.Duration {
	return config.QueryTimeout
}

// SetQueryTimeout allows runtime configuration of query timeout
func SetQueryTimeout(d time.Duration) {
	config.QueryTimeout = d
}

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

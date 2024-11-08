package directive

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/light-speak/lighthouse/cache"
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
)

// buildCacheKey builds a cache key from context, field name and args
func buildCacheKey(ctx *context.Context, field *ast.Field, isAuth bool, args string) string {
	var cacheKey strings.Builder
	cacheKey.Grow(512)

	if ctx.OprationName != nil {
		cacheKey.WriteString(*ctx.OprationName)
		cacheKey.WriteByte(':')
	}
	if isAuth && ctx.UserId != nil {
		cacheKey.WriteString(fmt.Sprint(*ctx.UserId))
		cacheKey.WriteByte(':')
	}
	cacheKey.WriteString(field.Name)
	cacheKey.WriteByte(':')
	cacheKey.WriteString(args)

	return cacheKey.String()
}

// buildArgsKey builds a sorted key from field arguments
func buildArgsKey(args map[string]*ast.Argument) string {
	if len(args) == 0 {
		return ""
	}

	argPairs := make([]string, 0, len(args))
	for _, arg := range args {
		argPairs = append(argPairs, fmt.Sprintf("%v:%v", arg.Name, arg.Value))
	}
	sort.Strings(argPairs)
	return strings.Join(argPairs, ",")
}

func CacheAfter(ctx *context.Context, field *ast.Field, directive *ast.Directive, store *ast.NodeStore, parent ast.Node, result interface{}) errors.GraphqlErrorInterface {
	// Get TTL with default optimization
	ttl := 30 * time.Minute
	if ttlArg := directive.GetArg("ttl"); ttlArg != nil {
		ttl = time.Duration(ttlArg.Value.(int64)) * time.Second
	}

	// Get auth flag with single check
	isAuth := false
	if auth := directive.GetArg("auth"); auth != nil && auth.Value.(bool) && ctx.UserId != nil {
		isAuth = true
	}

	// Build args and cache keys
	argsKey := buildArgsKey(field.Args)
	cacheKey := buildCacheKey(ctx, field, isAuth, argsKey)

	// Check if cache already exists
	var existingCache interface{}
	if err := cache.Get(cacheKey, &existingCache); err == nil {
		// Cache already exists, skip setting
		log.Info().
			Str("cacheKey", cacheKey).
			Str("field", field.Name).
			Bool("isAuth", isAuth).
			Msg("Cache already exists, skipping set")
		return nil
	}

	// Process tags with capacity pre-allocation
	tags := directive.GetArg("tags").Value.([]interface{})
	cacheTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagStr := tag.(string)
		if isAuth {
			tagStr = fmt.Sprintf("%v:%v", *ctx.UserId, tagStr)
		}
		cacheTags = append(cacheTags, tagStr)
	}

	// Set cache with retry mechanism
	maxRetries := 3
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = cache.Set(cacheKey, result, ttl, cacheTags...); err == nil {
			break
		}
		if i == maxRetries-1 {
			log.Error().Err(err).
				Str("field", field.Name).
				Str("attempt", fmt.Sprintf("%d/%d", i+1, maxRetries)).
				Msg("Failed to set cache after retries")
			return nil
		}
		time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
	}

	log.Info().
		Str("cacheKey", cacheKey).
		Strs("tags", cacheTags).
		Dur("ttl", ttl).
		Str("field", field.Name).
		Bool("isAuth", isAuth).
		Msg("Cache set successfully")

	return nil
}

func CacheBefore(ctx *context.Context, field *ast.Field, directive *ast.Directive, store *ast.NodeStore, parent ast.Node, resultChan interface{}) errors.GraphqlErrorInterface {
	// Get auth flag
	isAuth := false
	if auth := directive.GetArg("auth"); auth != nil && auth.Value.(bool) && ctx.UserId != nil {
		isAuth = true
	}

	// Build args and cache keys
	argsKey := buildArgsKey(field.Args)
	cacheKey := buildCacheKey(ctx, field, isAuth, argsKey)

	// Try to get from cache
	var cachedResult interface{}
	if err := cache.Get(cacheKey, &cachedResult); err == nil {
		// Cache hit, send result through channel
		resultChannel := resultChan.(chan interface{})
		resultChannel <- cachedResult
		log.Info().
			Str("cacheKey", cacheKey).
			Str("field", field.Name).
			Bool("isAuth", isAuth).
			Msg("Cache hit")
	}

	return nil
}

func init() {
	ast.AddFieldRuntimeBeforeDirective("cache", CacheBefore)
	ast.AddFieldRuntimeAfterDirective("cache", CacheAfter)
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/log"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	isEnabled   bool
)

func init() {
	config := env.LighthouseConfig.Redis
	if config.Host == "" {
		log.Warn().Msg("Redis host not configured, cache disabled")
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to connect to Redis, cache disabled")
		return
	}

	isEnabled = true
	log.Info().Msg("Redis cache enabled")
}

// IsEnabled returns whether cache is available
func IsEnabled() bool {
	return isEnabled
}

// Set stores value in cache with expiration and tags
func Set(key string, value interface{}, expiration time.Duration, tags ...string) error {
	if !isEnabled {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	ctx := context.Background()
	pipe := redisClient.Pipeline()

	// 设置缓存值
	pipe.Set(ctx, key, data, expiration)

	// 为每个tag添加关联的key
	for _, tag := range tags {
		pipe.SAdd(ctx, fmt.Sprintf("tag:%s", tag), key)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// Get retrieves value from cache
func Get(key string, result interface{}) error {
	if !isEnabled {
		return fmt.Errorf("cache is disabled")
	}

	ctx := context.Background()
	data, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache key not found: %s", key)
		}
		return err
	}

	return json.Unmarshal(data, result)
}

// Delete removes key from cache
func Delete(key string) error {
	if !isEnabled {
		return nil
	}

	ctx := context.Background()
	return redisClient.Del(ctx, key).Err()
}

// Clear removes all keys from cache
func Clear() error {
	if !isEnabled {
		return nil
	}

	ctx := context.Background()
	return redisClient.FlushAll(ctx).Err()
}

// DeleteByTags removes all keys associated with given tags
func DeleteByTags(tags ...string) error {
	if !isEnabled {
		return nil
	}

	ctx := context.Background()
	pipe := redisClient.Pipeline()

	// 收集所有要删除的键
	var keysToDelete []string
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s", tag)
		// 获取tag关联的所有键
		keys, err := redisClient.SMembers(ctx, tagKey).Result()
		if err != nil {
			return err
		}
		keysToDelete = append(keysToDelete, keys...)

		// 删除tag集合
		pipe.Del(ctx, tagKey)
	}

	// 删除所有关联的键
	if len(keysToDelete) > 0 {
		pipe.Del(ctx, keysToDelete...)
	}

	_, err := pipe.Exec(ctx)
	return err
}

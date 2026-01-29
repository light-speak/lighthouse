package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/logs"
	goRedis "github.com/redis/go-redis/v9"
)

type LightRedis struct {
	Client   *goRedis.Client
	IsEnable bool
}

var (
	LightRedisClient *LightRedis
	healthCheckStop  chan struct{}
	reconnectMu      sync.Mutex
)

func initRedis() {
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     LightRedisConfig.Host + ":" + LightRedisConfig.Port,
		Password: LightRedisConfig.Password,
		DB:       LightRedisConfig.DB,

		PoolSize:     LightRedisConfig.PoolSize,
		MinIdleConns: LightRedisConfig.MinIdleConns,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logs.Error().Err(err).Msg("failed to connect redis")
		LightRedisClient = &LightRedis{
			Client:   nil,
			IsEnable: false,
		}
		return
	}
	logs.Info().Msg("redis connected")
	LightRedisClient = &LightRedis{
		Client:   client,
		IsEnable: true,
	}

	// Start health check goroutine
	startRedisHealthCheck()
}

// startRedisHealthCheck starts a goroutine that monitors Redis health
// and attempts reconnection on consecutive failures
func startRedisHealthCheck() {
	if !LightRedisConfig.Enable {
		return
	}

	healthCheckStop = make(chan struct{})
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		defer ticker.Stop()

		failCount := 0
		for {
			select {
			case <-healthCheckStop:
				logs.Debug().Msg("redis health check stopped")
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := LightRedisClient.Client.Ping(ctx).Result()
				cancel()

				if err != nil {
					failCount++
					logs.Warn().Err(err).Int("fail_count", failCount).Msg("Redis health check failed")
					if failCount >= 3 {
						go reconnectRedis()
						failCount = 0
					}
				} else {
					if failCount > 0 {
						logs.Info().Int("previous_failures", failCount).Msg("Redis health recovered")
					}
					failCount = 0
				}
			}
		}
	}()
}

// reconnectRedis attempts to reconnect to Redis
func reconnectRedis() {
	reconnectMu.Lock()
	defer reconnectMu.Unlock()

	logs.Info().Msg("attempting Redis reconnection...")

	// Close old connection if exists
	if LightRedisClient != nil && LightRedisClient.Client != nil {
		_ = LightRedisClient.Client.Close()
	}

	// Create new client
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     LightRedisConfig.Host + ":" + LightRedisConfig.Port,
		Password: LightRedisConfig.Password,
		DB:       LightRedisConfig.DB,

		PoolSize:     LightRedisConfig.PoolSize,
		MinIdleConns: LightRedisConfig.MinIdleConns,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logs.Error().Err(err).Msg("Redis reconnection failed")
		return
	}

	LightRedisClient.Client = client
	LightRedisClient.IsEnable = true
	logs.Info().Msg("Redis reconnected successfully")
}

// StopHealthCheck stops the Redis health check goroutine
func StopHealthCheck() {
	if healthCheckStop != nil {
		close(healthCheckStop)
		healthCheckStop = nil
	}
}

func GetLightRedis() (*LightRedis, error) {
	if !LightRedisConfig.Enable {
		return nil, errors.New("redis is not enabled")
	}
	if LightRedisClient == nil {
		return nil, errors.New("redis is not initialized")
	}
	if !LightRedisClient.IsEnable {
		return nil, errors.New("redis is not enabled")
	}
	return LightRedisClient, nil
}

func GetClient() (*goRedis.Client, error) {

	if LightRedisClient == nil {
		return nil, errors.New("redis is not initialized")
	}
	if !LightRedisClient.IsEnable {
		return nil, errors.New("redis is not enabled")
	}
	return LightRedisClient.Client, nil
}

func (lr *LightRedis) GetClient() (*goRedis.Client, error) {
	if !LightRedisConfig.Enable {
		return nil, errors.New("redis is not enabled")
	}
	if LightRedisClient == nil {
		return nil, errors.New("redis is not initialized")
	}
	if !lr.IsEnable {
		return nil, errors.New("redis is not enabled")
	}
	return lr.Client, nil
}

func (lr *LightRedis) Close() error {
	// Stop health check first
	StopHealthCheck()

	if lr == nil || lr.Client == nil {
		return nil
	}
	return lr.Client.Close()
}

// Stats returns the Redis connection pool statistics
func (lr *LightRedis) Stats() map[string]any {
	if lr == nil || lr.Client == nil {
		return nil
	}
	stats := lr.Client.PoolStats()
	return map[string]any{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}
}

func (lr *LightRedis) Enable() bool {
	return lr.IsEnable
}

func (lr *LightRedis) Get(ctx context.Context, key string) (string, error) {
	client, err := lr.GetClient()
	if err != nil {
		return "", err
	}

	val, err := client.Get(ctx, key).Result()
	if err == goRedis.Nil {
		// key 不存在，不算错误
		return "", nil
	}
	return val, err
}

func (lr *LightRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, expiration).Err()
}

// nilMarker is used to cache nil values to prevent cache penetration
const nilMarker = "__NIL__"

// nilCacheTTL is the TTL for nil value cache (shorter than normal cache)
const nilCacheTTL = 5 * time.Minute

func Remember[T any](ctx context.Context, key string, callback func() T, expiration time.Duration) (*T, error) {
	redisClient, err := GetLightRedis()
	if err != nil {
		return nil, err
	}

	val, err := redisClient.Get(ctx, key)
	if err != nil {
		return nil, err // 这里已经过滤掉 redis.Nil
	}

	// Check if it's a cached nil value
	if val == nilMarker {
		return nil, nil
	}

	if val != "" {
		var result T
		if err := sonic.Unmarshal([]byte(val), &result); err == nil {
			return &result, nil
		}
		// 反序列化失败，清理坏缓存
		logs.Error().Err(err).Msg("failed to unmarshal value")
		_ = redisClient.Delete(ctx, key)
	}

	// 缓存未命中或坏数据 → 调用回调
	data := callback()

	if any(data) == nil {
		// Cache nil value with shorter TTL to prevent cache penetration
		ttl := nilCacheTTL
		if expiration > 0 && expiration < nilCacheTTL {
			ttl = expiration
		}
		_ = redisClient.Set(ctx, key, nilMarker, ttl)
		return nil, nil
	}

	bytes, err := sonic.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := redisClient.Set(ctx, key, bytes, expiration); err != nil {
		return nil, err
	}

	return &data, nil
}

func RememberForever[T any](ctx context.Context, key string, callback func() T) (*T, error) {
	return Remember(ctx, key, callback, 0)
}

// Delete removes key from redis
func (lr *LightRedis) Delete(ctx context.Context, key string) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.Del(ctx, key).Err()
}

// Clear clears all keys in current DB
func (lr *LightRedis) Clear(ctx context.Context) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.FlushDB(ctx).Err()
}

// Has checks if key exists
func (lr *LightRedis) Has(ctx context.Context, key string) bool {
	client, err := lr.GetClient()
	if err != nil {
		return false
	}
	val, err := client.Exists(ctx, key).Result()
	return err == nil && val > 0
}

// Increment increments value by 1
func (lr *LightRedis) Increment(ctx context.Context, key string) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.Incr(ctx, key).Err()
}

// IncrementBy increments value by given amount
func (lr *LightRedis) IncrementBy(ctx context.Context, key string, value int64) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.IncrBy(ctx, key, value).Err()
}

// Decrement decrements value by 1
func (lr *LightRedis) Decrement(ctx context.Context, key string) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.Decr(ctx, key).Err()
}

// DecrementBy decrements value by given amount
func (lr *LightRedis) DecrementBy(ctx context.Context, key string, value int64) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.DecrBy(ctx, key, value).Err()
}

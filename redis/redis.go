package redis

import (
	"context"
	"errors"
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/logs"
	goRedis "github.com/redis/go-redis/v9"
)

type LightRedis struct {
	Client   *goRedis.Client
	IsEnable bool
}

var LightRedisClient *LightRedis

func initRedis() {
	client := goRedis.NewClient(&goRedis.Options{
		Addr:     LightRedisConfig.Host + ":" + LightRedisConfig.Port,
		Password: LightRedisConfig.Password,
		DB:       LightRedisConfig.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logs.Error().Err(err).Msg("failed to connect redis")
	}
	logs.Info().Msg("redis connected")
	LightRedisClient = &LightRedis{
		Client:   client,
		IsEnable: true,
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

func (lr *LightRedis) Enable() bool {
	return lr.IsEnable
}

func (lr *LightRedis) Get(ctx context.Context, key string) (string, error) {
	client, err := lr.GetClient()
	if err != nil {
		return "", err
	}
	return client.Get(ctx, key).Result()
}

func (lr *LightRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	client, err := lr.GetClient()
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, expiration).Err()
}

// Remember gets cached value or stores the result of callback
func (lr *LightRedis) Remember(ctx context.Context, key string, callback func() interface{}, expiration time.Duration) (interface{}, error) {
	// Try to get from cache first
	val, err := lr.Get(ctx, key)
	if err == nil {
		var result interface{}
		err = sonic.Unmarshal([]byte(val), &result)
		if err == nil {
			return result, nil
		}
	}

	// Get fresh data
	data := callback()

	// Cache the result
	bytes, err := sonic.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = lr.Set(ctx, key, string(bytes), expiration)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// RememberForever caches value forever
func (lr *LightRedis) RememberForever(ctx context.Context, key string, callback func() interface{}) (interface{}, error) {
	return lr.Remember(ctx, key, callback, 0)
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

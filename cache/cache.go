package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Client *redis.Client
}

func NewCache() (*Cache, error) {
	opts, err := redis.ParseURL(fmt.Sprintf(
		"redis://%s:%s@%s:%s/%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PASS"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_DB"),
	))
	if err != nil {
		return nil, err
	}

	return &Cache{
		Client: redis.NewClient(opts),
	}, nil
}

func (c *Cache) CacheWeatherData(ctx context.Context, lat, long, data string) error {
	// Timeout the context after 5 seconds to avoid long blocking operations.
	// This is useful if the Redis server is down or unreachable.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Data is in JSON format.
	key := genCacheKey(lat, long)
	err := c.Client.SetEX(ctx, key, data, time.Hour).Err()
	return err
}

/// If the return value is ("", nil), it means the key did not have an associated value stored in cache.
func (c *Cache) GetCachedWeatherData(ctx context.Context, lat, long string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Return JSON string that is stored in redis.
	key := genCacheKey(lat, long)
	// If the error is something other than "key does not exist in the cache".
	json, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist in cache.
	} else if err != nil {
		return "", err // Some other error occurred.
	}

	return json, nil 
}

func (c *Cache) Close() error {
	return c.Client.Close()
}

func genCacheKey(lat, long string) string {
	return fmt.Sprintf("weather:%s:%s", lat, long)
}
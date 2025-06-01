package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func Init() error {
	opt, err := redis.ParseURL(fmt.Sprintf(
		"redis://%s:%s@%s:%s/%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PASS"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_DB"),
	))

	if err != nil {
		return err
	}

	client = redis.NewClient(opt)

	return nil
}

func CacheWeatherData(ctx context.Context, lat, long, data string) error {
	// Data is in JSON format.
	key := lat + long
	err := client.SetEX(ctx, key, data, 12 * time.Hour).Err()
	return err
}

/// If the return value is ("", nil), it means the key did not have an associated value stored in cache.
func GetCachedWeatherData(ctx context.Context, lat, long string) (string, error) {
	// Return JSON string that is stored in redis.
	key := lat + long
	var (
		json string
		err error
	)
	// If the error is something other than "key does not exist in the cache".
	if json, err = client.Get(ctx, key).Result(); err != nil && err != redis.Nil {
		return "", err
	} else {
		return json, nil
	}
}
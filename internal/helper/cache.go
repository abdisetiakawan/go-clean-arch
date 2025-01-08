package helper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheHelper struct {
	client *redis.Client
}

func NewCacheHelper(client *redis.Client) *CacheHelper {
	return &CacheHelper{client: client}
}

func (c *CacheHelper) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *CacheHelper) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *CacheHelper) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *CacheHelper) GetAndUnmarshal(ctx context.Context, key string, value interface{}) error {
	cachedData, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cachedData), value)
}
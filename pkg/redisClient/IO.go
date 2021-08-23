package redisConnection

import (
	"context"
	"time"
)

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.rdb.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return rc.rdb.Get(ctx, key).Bytes()
}

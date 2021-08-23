package redisConnection

import "github.com/go-redis/redis/v8"

type RedisClient struct {
	rdb *redis.Client
}

func Connect(options *RedisOptions) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return &RedisClient{rdb: rdb}
}

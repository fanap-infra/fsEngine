package redisClient

import "github.com/go-redis/redis/v8"

func Connect(options *RedisOptions) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	return rdb
}

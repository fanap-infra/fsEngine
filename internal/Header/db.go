package Header_

import (
	"context"
	"time"
)

type RedisDB interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}

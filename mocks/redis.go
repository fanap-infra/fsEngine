package mocks

import (
	"context"
	"fmt"
	"time"
)

type RedisMock struct {
	dataStore map[string][]byte
}

func NewRedisMock() RedisMock {
	return RedisMock{dataStore: make(map[string][]byte)}
}

func (redisMock *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	redisMock.dataStore[key] = value.([]byte)
	return nil
}

func (redisMock *RedisMock) Get(ctx context.Context, key string) ([]byte, error) {
	value, isExist := redisMock.dataStore[key]
	if !isExist {
		return nil, fmt.Errorf("there is no value")
	}

	return value, nil
}

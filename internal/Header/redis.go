package Header_

import (
	"context"
)

func (hfs *HFileSystem) setRedisKeyValue(key string, value []byte) error {
	ctx := context.Background()

	err := hfs.redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		hfs.log.Errorv("Can not Set value in redis", "err", err.Error())
		return err
	}

	return nil
}

func (hfs *HFileSystem) getRedisValue(key string) ([]byte, error) {
	ctx := context.Background()

	data, err := hfs.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		hfs.log.Errorv("Can not get value in redis", "err", err.Error())
		return nil, err
	}

	return data, nil
}

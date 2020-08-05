package storage

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type objectCache struct {
	ctx context.Context
}

func newCache() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": "192.168.0.102:6379",
		},
	})

	c := cache.New(&cache.Options{
		Redis: ring,
	})

	return c
}

func (o *objectCache) setCacheObject(minioObject []byte, objectName string) error {
	err := newCache().Set(&cache.Item{
		Ctx:   o.ctx,
		Key:   objectName,
		Value: minioObject,
		TTL:   time.Hour,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *objectCache) getCacheObject(objectName string) ([]byte, error) {
	var buf []byte
	err := newCache().Get(o.ctx, objectName, &buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (o *objectCache) deleteCacheObject(objectName string) error {
	err := newCache().Delete(o.ctx, objectName)
	if err != nil {
		return err
	}

	return nil
}

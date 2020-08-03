package storage

import (
	"context"
	"log"
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
			"server": ":6379",
		},
	})

	c := cache.New(&cache.Options{
		Redis: ring,
	})

	return c
}

func (o *objectCache) setCacheObject(minioObject []byte, objectName string) {
	err := newCache().Set(&cache.Item{
		Ctx:   o.ctx,
		Key:   objectName,
		Value: minioObject,
		TTL:   time.Hour,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func (o *objectCache) getCacheObject(objectName string) []byte {
	var buf []byte
	err := newCache().Get(o.ctx, objectName, &buf)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return buf
}

func (o *objectCache) deleteCacheObject(objectName string) {
	err := newCache().Delete(o.ctx, objectName)
	if err != nil {
		log.Println(err)
	}
}

package bucket

import (
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cacher interface {
	PutObject(objectName string, minioObject []byte) error
	GetObject(objectName string) ([]byte, error)
	RemoveObject(objectName string) error
	InitCache(CacheConfig) error
}

type CacheConfig struct {
	Config map[string]interface{}
}

type cacheObject struct {
	ctx context.Context
	cache *cache.Cache
}

func (co *cacheObject) InitCache() error {
	var b Bucket
	ac, err := b.GetConfig()
	if err != nil {
		return err
	}

	endpoint := ac.config["cacheServerEndpoint"]
	if endpoint == "" {
		return errors.New("new cache failed")
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": endpoint.(string),
		},
	})

	c := cache.New(&cache.Options{
		Redis: ring,
	})

	co.cache = c
	return nil
}

func (co *cacheObject) PutObject(objectName string, minioObject []byte) error {
	err := co.cache.Set(&cache.Item{
		Ctx:   co.ctx,
		Key:   objectName,
		Value: minioObject,
		TTL:   time.Hour,
	})
	if err != nil {
		return err
	}
	return nil
}

func (co *cacheObject) GetObject(objectName string) ([]byte, error) {
	var buf []byte
	err := co.cache.Get(co.ctx, objectName, &buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (co *cacheObject) RemoveObject(objectName string) error {
	err := co.cache.Delete(co.ctx, objectName)
	if err != nil {
		return err
	}

	return nil
}
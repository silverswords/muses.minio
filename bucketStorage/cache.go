package bucketStorage

import (
	"context"
	"errors"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Cacher interface {
	PutObject(objectName string, minioObject []byte) error
	GetObject(objectName string) ([]byte, error)
	RemoveObject(objectName string) error
}

type cacheObject struct {
	cache *cache.Cache
}

func initCache(configName, configPath string) (Cacher, error) {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return nil, err
	}
	cacheType := ac.Cache["cacheType"]

	if cacheType.(string) == "redis" {
		c, err := newRedisCache(configName, configPath)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return nil, nil
}

func newRedisCache(configName, configPath string) (Cacher, error) {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return nil, err
	}

	endpoint := ac.Cache["cacheServerEndpoint"]
	if endpoint == "" {
		return nil, errors.New("new cache failed")
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": endpoint.(string),
		},
	})

	c := cache.New(&cache.Options{
		Redis: ring,
	})

	return &cacheObject{
		cache: c,
	}, nil
}

func (co *cacheObject) PutObject(objectName string, minioObject []byte) error {
	log.Println("minioObject", minioObject)
	err := co.cache.Set(&cache.Item{
		Ctx:   context.Background(),
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
	err := co.cache.Get(context.Background(), objectName, &buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (co *cacheObject) RemoveObject(objectName string) error {
	err := co.cache.Delete(context.Background(), objectName)
	if err != nil {
		return err
	}

	return nil
}
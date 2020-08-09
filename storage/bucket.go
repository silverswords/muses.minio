package storage

import (
	"context"
	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName string
	objectCache
	configInfo
	OtherOptions
}

type OtherOptions struct {
	location string
	strategy string
	cache    bool
}

type OtherOption func(*OtherOptions)

func WithStrategy(strategy string) OtherOption {
	return func(b *OtherOptions) {
		b.strategy = strategy
	}
}

func WithLocation(location string) OtherOption {
	return func(b *OtherOptions) {
		b.location = location
	}
}

func WithCache(cache bool) OtherOption {
	return func(b *OtherOptions) {
		b.cache = cache
	}
}

func NewBucketConfig(bucketName, configName, configPath string, opts ...OtherOption) *Bucket {
	const(
		defaultStrategy = "multiWriteStrategy"
		defaultLocation = "cn-north-1"
		defaultCache = false
	)

	b := &OtherOptions{
		defaultLocation,
		defaultStrategy,
		defaultCache,
	}

	for _, opt := range opts {
		opt(b)
	}

	ctx := context.TODO()

	return &Bucket{
		bucketName: bucketName,
		objectCache: objectCache{
			ctx,
		},
		configInfo: configInfo{
			configName,
			configPath,
		},
		OtherOptions: OtherOptions{
			b.location,
			b.strategy,
			b.cache,
		},
	}
}

func (b *Bucket) MakeBucket() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.MakeBucket(b.bucketName, b.location)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bucket) CheckBucket(bucketName string) (bool, error) {
	var exists bool
	clients, err := b.getStrategyClients()
	if err != nil {
		return false, err
	}

	for _, v := range clients {
		exists, err = v.client.BucketExists(bucketName)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}

	return exists, nil
}

func (b *Bucket) ListedBucket() ([]minio.BucketInfo, error) {
	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	minioClient := clients[0].client
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		return nil, err
	}

	return buckets, err
}

func (b *Bucket) RemoveBucket(bucketName string) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.RemoveBucket(bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

package storage

import (
	"context"
	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName string
	location   string
	objectCache
	strategy string
	configInfo
}

func NewBucket(bucketName, location, strategy, configName, configPath string) *Bucket {
	ctx := context.TODO()
	return &Bucket{
		bucketName: bucketName,
		location:   location,
		objectCache: objectCache{
			ctx,
		},
		strategy: strategy,
		configInfo: configInfo{
			configName: configName,
			configPath: configPath,
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
	var err error
	clients, err := b.getStrategyClients()
	if err != nil {
		return false, err
	}

	for _, v := range clients {
		exists, err = v.client.BucketExists(bucketName)
		if err != nil {
			return false, err
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

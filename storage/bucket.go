package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

type Bucket struct {
	bucketName string
	objectCache
	clientConfigInfo
	OtherBucketConfigOptions
}

type OtherBucketConfigOptions struct {
	region           string
	Strategy         string
	cache            bool
	bucketEncryption bool
	expirationDays   lifecycle.ExpirationDays
}

func NewBucketConfig(bucketName, configName, configPath string, opts OtherBucketConfigOptions) *Bucket {
	ctx := context.TODO()

	return &Bucket{
		bucketName: bucketName,
		objectCache: objectCache{
			ctx,
		},
		clientConfigInfo: clientConfigInfo{
			configName,
			configPath,
		},
		OtherBucketConfigOptions: opts,
	}
}

func (b *Bucket) MakeBucket() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.MakeBucket(context.Background(), b.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	if b.bucketEncryption {
		err := b.setBucketEncryption()
		if err != nil {
			return err
		}
	}

	err = b.setBucketLifecycle(b.expirationDays)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) CheckBucket() (bool, error) {
	var exists bool
	clients, err := b.getStrategyClients()
	if err != nil {
		return false, err
	}

	for _, v := range clients {
		exists, err = v.client.BucketExists(context.Background(), b.bucketName)
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
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	return buckets, err
}

func (b *Bucket) RemoveBucket() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.RemoveBucket(context.Background(), b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

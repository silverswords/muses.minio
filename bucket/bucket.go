package bucket

import (
	"context"
	"github.com/minio/minio-go/v7"
	"os"
)

type Bucket struct {
	client  client
	cache   cacher
	bucketOperator bucketOperator
	bucketName string
	clientConfigInfo
	minioClient
	cacheObject
	OtherBucketConfigOptions
}

type OtherBucketConfigOptions struct {
	region           string
	cache            bool
}

func NewBucketConfig(bucketName, configName, configPath string, opts OtherBucketConfigOptions) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		clientConfigInfo: clientConfigInfo{
			configName,
			configPath,
		},
		cacheObject: cacheObject{
			ctx: context.Background(),
		},
		OtherBucketConfigOptions: OtherBucketConfigOptions{
			opts.region,
			opts.cache,
		},
	}
}

func (b *Bucket) MakeBucket() error {
	err := b.bucketOperator.MakeBucket(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) CheckBucket() (bool, error) {
	exists, err := b.bucketOperator.CheckBucket(b.bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (b *Bucket) ListedBucket() ([]minio.BucketInfo, error) {
	bucketInfos, err := b.bucketOperator.ListBuckets()
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func (b *Bucket) RemoveBucket() error {
	err := b.bucketOperator.RemoveBucket(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) PutObject(objectName string, object *os.File) error {
	err := b.client.PutObject(b.bucketName, objectName, object)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	var buf []byte
	buf, err := b.client.GetObject(b.bucketName, objectName)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (b *Bucket) RemoveObject(objectName string) error {
	err := b.client.RemoveObject(b.bucketName, objectName)
	if err != nil {
		return err
	}
}

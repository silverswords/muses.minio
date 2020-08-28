package bucket

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/minio/minio-go/v7/pkg/replication"
	"io"
	"os"
)

type Bucket struct {
	client  Client
	cacher   Cacher
	bucketName string
	clientConfigInfo
	minioClient
	cacheObject
	OtherBucketConfigOptions
}

func NewBucketConfig(bucketName, configName, configPath string, opts ...OtherBucketConfigOption) (*Bucket, error) {
	const (
		defaultUseCache = false
	)

	var m minioClient
	err := m.InitClient()
	if err != nil {
		return nil, err
	}

	b := &OtherBucketConfigOptions{
		defaultUseCache,
	}

	for _, opt := range opts {
		opt(b)
	}

	var bc Bucket
	if b.useCache {
		bc.InitCache()
	}

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
			b.useCache,
		},
	}, nil
}

func (b *Bucket) MakeBucket(opts ...OtherMakeBucketOption) error {
	const (
		defaultRegion = "us-east-1"
		defaultObjectLocking = false
	)

	o := &OtherMakeBucketOptions{
		Region: defaultRegion,
		ObjectLocking: defaultObjectLocking,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.client.MakeBucket(b.bucketName, o)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) CheckBucket() (bool, error) {
	exists, err := b.client.CheckBucket(b.bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (b *Bucket) ListedBucket() ([]minio.BucketInfo, error) {
	bucketInfos, err := b.client.ListBuckets()
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func (b *Bucket) RemoveBucket() error {
	err := b.client.RemoveBucket(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) SetBucketReplication(cfg replication.Config) error {
	err := b.client.SetBucketReplication(b.bucketName, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketReplication() (replication.Config, error) {
	cfg, err := b.client.GetBucketReplication(b.bucketName)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (b *Bucket) RemoveBucketReplication() error {
	err := b.client.RemoveBucketReplication(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) PutObject(objectName string, object *os.File, opts ...OtherPutObjectOption) error {
	stat, err := object.Stat()
	buf := make([]byte, stat.Size())
	_, err = io.ReadFull(object, buf)
	if err != nil {
		return err
	}

	if b.useCache {
		err = b.cacher.PutObject(objectName, buf)
		if err != nil {
			return err
		}
	}

	var e encrypt.ServerSide
	var (
		defaultServerSideEncryption = e
	)

	o := &OtherPutObjectOptions{
		defaultServerSideEncryption,
	}

	for _, opt := range opts {
		opt(o)
	}

	err = b.client.PutObject(b.bucketName, objectName, object, o)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObject(objectName string, opts ...OtherGetObjectOption) ([]byte, error) {
	var buf []byte
	if b.useCache {
		buf, err := b.cacher.GetObject(objectName)
		if err != nil {
			return nil, err
		}

		if buf != nil {
			return buf, nil
		}
	}

	var e encrypt.ServerSide
	o := &OtherGetObjectOptions{
		e,
	}

	for _, opt := range opts {
		opt(o)
	}

	buf, err := b.client.GetObject(b.bucketName, objectName, o)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (b *Bucket) RemoveObject(objectName string, opts ...OtherRemoveObjectOption) error {
	const (
		defaultGovernanceBypass = false
	)

	o := &OtherRemoveObjectOptions{
		GovernanceBypass: defaultGovernanceBypass,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.client.RemoveObject(b.bucketName, objectName, o)
	if err != nil {
		return err
	}

	if b.useCache {
		err = b.cacher.RemoveObject(objectName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) ListObjects(bucketName string, opts ...OtherListObjectsOption) <-chan minio.ObjectInfo {
	const (
		defaultPrefix = ""
	)

	o := &OtherListObjectsOptions{
		prefix: defaultPrefix,
	}

	for _, opt := range opts {
		opt(o)
	}

	objectInfo := b.client.ListObjects(bucketName, o)

	return objectInfo
}

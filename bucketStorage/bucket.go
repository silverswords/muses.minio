package bucketStorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/silverswords/muses.minio/bucketStorage/driver"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Bucket struct {
	client     Client
	cacher     Cacher
	bucketName string
	ConfigInfo
	mc MinioClient
}

func NewBucket(bucketName, configName, configPath string) (*Bucket, error) {
	c, err := initClient(configName, configPath)
	if err != nil {
		return nil, err
	}

	ca, err := initCache(configName, configPath)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("cacher:", ca)
	return &Bucket{
		client:     c,
		cacher:     ca,
		bucketName: bucketName,
		ConfigInfo: ConfigInfo{
			configName,
			configPath,
		},
	}, nil
}

func (b *Bucket) MakeBucket(opts ...OtherMakeBucketOption) error {
	const (
		defaultRegion        = "us-east-1"
		defaultObjectLocking = false
	)

	o := &MakeBucketOptions{
		Region:        defaultRegion,
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
	exists, err := b.mc.CheckBucket(b.bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (b *Bucket) ListedBucket() ([]minio.BucketInfo, error) {
	bucketInfos, err := b.mc.ListBuckets()
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func (b *Bucket) RemoveBucket() error {
	err := b.mc.RemoveBucket(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) SetBucketVersioning(opts ...OtherSetBucketVersioningOption) error {
	const (
		defaultStatus = "Enabled"
	)

	o := &SetBucketVersioningOptions{
		Status: defaultStatus,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.mc.SetBucketVersioning(b.bucketName, o)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketVersioning() (minio.BucketVersioningConfiguration, error) {
	configuration, err := b.mc.GetBucketVersioning(b.bucketName)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func (b *Bucket) SetBucketReplication(cfg replication.Config) error {
	err := b.mc.SetBucketReplication(b.bucketName, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketReplication() (replication.Config, error) {
	cfg, err := b.mc.GetBucketReplication(b.bucketName)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (b *Bucket) RemoveBucketReplication() error {
	err := b.mc.RemoveBucketReplication(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) SetBucketPolicy(policy string) error {
	err := b.mc.SetBucketPolicy(b.bucketName, policy)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketPolicy() (string, error) {
	policy, err := b.mc.GetBucketPolicy(b.bucketName)
	if err != nil {
		return "", err
	}

	return policy, nil
}

func (b *Bucket) SetObjectLockConfig(mode string, validity *uint, uint string) error {
	err := b.mc.SetObjectLockConfig(b.bucketName, mode, validity, uint)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObjectLockConfig() (string, string, *uint, string, error) {
	objectLock, mode, validity, uint, err := b.mc.GetObjectLockConfig(b.bucketName)
	if err != nil {
		return "", "", nil, "", err
	}

	return objectLock, mode, validity, uint, nil
}

func (b *Bucket) NewTypedWriter(ctx context.Context, key string, opts ...OtherPutObjectOption) (driver.Writer, error) {
	var buf bytes.Buffer

	var e encrypt.ServerSide
	var d = int64(1024)
	var (
		defaultServerSideEncryption = e
		defaultObjectSize = d
	)

	o := &OtherPutObjectOptions{
		defaultServerSideEncryption,
		defaultObjectSize,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.client.PutObject(ctx, b.bucketName, key, &buf, o)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *Bucket) SignedURL(ctx context.Context, key string, expires time.Duration, Method string) (string, error) {
	var u *url.URL
	var err error
	switch Method {
	case http.MethodGet:
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\"file\"")
		u, err = b.mc.PresignedGetObject(ctx, b.bucketName, key, expires, reqParams)
		if err != nil {
			return "", err
		}
	case http.MethodPut:
		u, err = b.mc.PresignedPutObject(ctx, b.bucketName, key, expires)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported Method #{opts.Method}")
	}
	if err != nil {
		return "", err
	}

	return u.Path, nil
}

func (b *Bucket) NewRangeReader(ctx context.Context, key string, opts ...OtherGetObjectOption) (driver.Reader, error) {
	var e encrypt.ServerSide
	o := &GetObjectOptions{
		e,
	}

	for _, opt := range opts {
		opt(o)
	}

	object, err := b.client.GetObject(b.bucketName, key, o)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (b *Bucket) Delete(ctx context.Context, key string, opts ...OtherRemoveObjectOption) error {
	const (
		defaultGovernanceBypass = false
	)

	o := &RemoveObjectOptions{
		GovernanceBypass: defaultGovernanceBypass,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.client.RemoveObject(ctx, b.bucketName, key, o)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) ListObjects(opts ...OtherListObjectsOption) <-chan minio.ObjectInfo {
	const (
		defaultPrefix = ""
	)

	o := &ListObjectsOptions{
		prefix: defaultPrefix,
	}

	for _, opt := range opts {
		opt(o)
	}

	objectInfo := b.client.ListObjects(b.bucketName, o)

	return objectInfo
}

func (b *Bucket) Close() error {
	return nil
}
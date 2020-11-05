package bucketStorage

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/minio/minio-go/v7/pkg/replication"
	"io"
	"net/url"
	"time"
)

type Bucketer interface {
	Uploader
	Downloader
	Remover
}

type Uploader interface {
	PutObject(objectName string, reader io.Reader, opts ...OtherPutObjectOption) error
}

type Downloader interface {
	GetObject(objectName string, opts ...OtherGetObjectOption) ([]byte, error)
}

type Remover interface {
	RemoveObject(objectName string, opts ...OtherRemoveObjectOption) error
}

type Bucket struct {
	client     Client
	cacher     Cacher
	bucketName string
	ConfigInfo
	minioClient
}

func NewBucket(bucketName, configName, configPath string) (*Bucket, error) {
	c, err := initClient(configName, configPath)
	if err != nil {
		return nil, err
	}

	ca, err := initCache(configName, configPath)
	if err != nil {
		return nil, err
	}

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

	err := b.client.SetBucketVersioning(b.bucketName, o)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketVersioning() (minio.BucketVersioningConfiguration, error) {
	configuration, err := b.client.GetBucketVersioning(b.bucketName)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
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

func (b *Bucket) SetBucketPolicy(policy string) error {
	err := b.client.SetBucketPolicy(b.bucketName, policy)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetBucketPolicy() (string, error) {
	policy, err := b.client.GetBucketPolicy(b.bucketName)
	if err != nil {
		return "", err
	}

	return policy, nil
}

func (b *Bucket) SetObjectLockConfig(mode string, validity *uint, uint string) error {
	err := b.client.SetObjectLockConfig(b.bucketName, mode, validity, uint)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObjectLockConfig() (string, string, *uint, string, error) {
	objectLock, mode, validity, uint, err := b.client.GetObjectLockConfig(b.bucketName)
	if err != nil {
		return "", "", nil, "", err
	}

	return objectLock, mode, validity, uint, nil
}

func (b *Bucket) PutObject(objectName string, reader io.Reader, opts ...OtherPutObjectOption) error {
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

	cacheBytes := make([]byte, o.ObjectSize)

	teeReader := io.TeeReader(reader, &buf)
	_, err := teeReader.Read(cacheBytes)
	if err != nil {
		return err
	}

	err = b.client.PutObject(b.bucketName, objectName, &buf, o)
	if err != nil {
		return err
	}

	err = b.cacher.PutObject(objectName, cacheBytes)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) PresignedPutObject(ctx context.Context, objectName string, expires time.Duration) (*url.URL, error) {
	url, err := b.client.PresignedPutObject(b.bucketName, objectName, expires)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (b *Bucket) PresignedGetObject(ctx context.Context, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	url, err := b.client.PresignedGetObject(b.bucketName, objectName, expires, reqParams)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (b *Bucket) GetObject(objectName string, opts ...OtherGetObjectOption) ([]byte, error) {
	var buf []byte
	var e encrypt.ServerSide
	o := &GetObjectOptions{
		e,
	}

	for _, opt := range opts {
		opt(o)
	}

	buf, err := b.cacher.GetObject(objectName)
	if err != nil {
		return nil, err
	}
	if buf != nil {
		return buf, nil
	}

	buf, err = b.client.GetObject(b.bucketName, objectName, o)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (b *Bucket) RemoveObject(objectName string, opts ...OtherRemoveObjectOption) error {
	const (
		defaultGovernanceBypass = false
	)

	o := &RemoveObjectOptions{
		GovernanceBypass: defaultGovernanceBypass,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.client.RemoveObject(b.bucketName, objectName, o)
	if err != nil {
		return err
	}

	err = b.cacher.RemoveObject(objectName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) ListObjects(bucketName string, opts ...OtherListObjectsOption) <-chan minio.ObjectInfo {
	const (
		defaultPrefix = ""
	)

	o := &ListObjectsOptions{
		prefix: defaultPrefix,
	}

	for _, opt := range opts {
		opt(o)
	}

	objectInfo := b.client.ListObjects(bucketName, o)

	return objectInfo
}

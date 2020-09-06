package bucketStorage

import (
	"bytes"
	"crypto/md5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/minio/minio-go/v7/pkg/replication"
	"io"
)

type CacheBucket struct {
	cacher Cacher
	Bucket
}

type Bucket struct {
	client  Client
	bucketName string
	ConfigInfo
	minioClient
}

func NewBucketConfig(bucketName, configName, configPath string) (*Bucket, error) {
	c, err := initClient(configName, configPath)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		client: c,
		bucketName: bucketName,
		ConfigInfo: ConfigInfo{
			configName,
			configPath,
		},
	}, nil
}

func NewBucketWithCacheConfig(bucketName, configName, configPath string) (*CacheBucket, error) {
	c, err := initClient(configName, configPath)
	if err != nil {
		return nil, err
	}

	ca, err := initCache(configName, configPath)
	if err != nil {
		return nil, err
	}

	return &CacheBucket{
		cacher: ca,
		Bucket: Bucket{
			client: c,
			bucketName: bucketName,
			ConfigInfo: ConfigInfo{
				configName,
				configPath,
			},
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

func (b *Bucket) SetBucketVersioning(opts ...OtherSetBucketVersioningOption) error {
	const (
		defaultStatus = "Enabled"
	)

	o := &OtherSetBucketVersioningOptions{
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

func (b *Bucket) SetObjectLockConfig(mode *minio.RetentionMode, validity *uint, uint *minio.ValidityUnit) error {
	err := b.client.SetObjectLockConfig(b.bucketName, mode, validity, uint)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObjectLockConfig() (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error) {
	objectLock, mode, validity, uint, err := b.client.GetObjectLockConfig(b.bucketName)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return objectLock, mode, validity, uint, nil
}

func (cb *CacheBucket) PutObject(objectName string, reader io.Reader, objectSize int64, opts ...OtherPutObjectOption) error {
	var buf bytes.Buffer
	cacheBytes := make([]byte, objectSize)

	teeReader := io.TeeReader(reader, &buf)
	_, err := teeReader.Read(cacheBytes)
	if err != nil {
		return err
	}

	err = cb.Bucket.PutObject(objectName, &buf, objectSize, opts...)
	if err != nil {
		return err
	}

	err = cb.cacher.PutObject(objectName, cacheBytes)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) PutObject(objectName string, reader io.Reader, objectSize int64, opts ...OtherPutObjectOption) error {
	var buf bytes.Buffer
	cacheBytes := make([]byte, objectSize)
	teeReader := io.TeeReader(reader, &buf)
	_, err := teeReader.Read(cacheBytes)
	if err != nil {
		return err
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

	err = b.client.PutObject(b.bucketName, objectName, &buf, objectSize,  o)
	if err != nil {
		return err
	}

	h := md5.New()
	_, err = io.WriteString(h, objectName)
	if err != nil {
		return err
	}

	_, err = h.Write(cacheBytes)
	if err != nil {
		return err
	}

	md5.Sum(nil)

	return nil
}

func (cb *CacheBucket) GetObject(objectName string, opts ...OtherGetObjectOption) ([]byte, error) {
	var buf []byte
	buf, err := cb.cacher.GetObject(objectName)
	if err != nil {
		return nil, err
	}
	if buf != nil {
		return buf, nil
	}

	buf, err = cb.Bucket.GetObject(objectName, opts...)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (b *Bucket) GetObject(objectName string, opts ...OtherGetObjectOption) ([]byte, error) {
	var buf []byte
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

func (cb *CacheBucket) RemoveObject(objectName string, opts ...OtherRemoveObjectOption) error {
	err := cb.Bucket.RemoveObject(objectName, opts...)
	if err != nil {
		return err
	}

	err = cb.cacher.RemoveObject(objectName)
	if err != nil {
		return err
	}

	return nil
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

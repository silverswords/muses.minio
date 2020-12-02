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

type Storage interface {
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
	MinioClient
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

func MakeBucket(b *Bucket, opts ...OtherMakeBucketOption) error {
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

func CheckBucket(b *Bucket) (bool, error) {
	exists, err := b.CheckBucket(b.bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func ListedBucket(b *Bucket) ([]minio.BucketInfo, error) {
	bucketInfos, err := b.ListBuckets()
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func RemoveBucket(b *Bucket) error {
	err := b.RemoveBucket(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func SetBucketVersioning(b *Bucket, opts ...OtherSetBucketVersioningOption) error {
	const (
		defaultStatus = "Enabled"
	)

	o := &SetBucketVersioningOptions{
		Status: defaultStatus,
	}

	for _, opt := range opts {
		opt(o)
	}

	err := b.SetBucketVersioning(b.bucketName, o)
	if err != nil {
		return err
	}

	return nil
}

func GetBucketVersioning(b *Bucket) (minio.BucketVersioningConfiguration, error) {
	configuration, err := b.GetBucketVersioning(b.bucketName)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func SetBucketReplication(b *Bucket, cfg replication.Config) error {
	err := b.SetBucketReplication(b.bucketName, cfg)
	if err != nil {
		return err
	}

	return nil
}

func GetBucketReplication(b *Bucket) (replication.Config, error) {
	cfg, err := b.GetBucketReplication(b.bucketName)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func RemoveBucketReplication(b *Bucket) error {
	err := b.RemoveBucketReplication(b.bucketName)
	if err != nil {
		return err
	}

	return nil
}

func SetBucketPolicy(b *Bucket, policy string) error {
	err := b.SetBucketPolicy(b.bucketName, policy)
	if err != nil {
		return err
	}

	return nil
}

func GetBucketPolicy(b *Bucket) (string, error) {
	policy, err := b.GetBucketPolicy(b.bucketName)
	if err != nil {
		return "", err
	}

	return policy, nil
}

func SetObjectLockConfig(b *Bucket, mode string, validity *uint, uint string) error {
	err := b.SetObjectLockConfig(b.bucketName, mode, validity, uint)
	if err != nil {
		return err
	}

	return nil
}

func GetObjectLockConfig(b *Bucket) (string, string, *uint, string, error) {
	objectLock, mode, validity, uint, err := b.GetObjectLockConfig(b.bucketName)
	if err != nil {
		return "", "", nil, "", err
	}

	return objectLock, mode, validity, uint, nil
}

func PutObject(b *Bucket, objectName string, reader io.Reader, opts ...OtherPutObjectOption) error {
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

func PresignedPutObject(b *Bucket, ctx context.Context, objectName string, expires time.Duration) (*url.URL, error) {
	url, err := b.PresignedPutObject(b.bucketName, objectName, expires)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func PresignedGetObject(b *Bucket, ctx context.Context, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	url, err := b.PresignedGetObject(b.bucketName, objectName, expires, reqParams)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func GetObject(b *Bucket, objectName string, opts ...OtherGetObjectOption) ([]byte, error) {
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

func RemoveObject(b *Bucket, objectName string, opts ...OtherRemoveObjectOption) error {
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

func ListObjects(b *Bucket, opts ...OtherListObjectsOption) <-chan minio.ObjectInfo {
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

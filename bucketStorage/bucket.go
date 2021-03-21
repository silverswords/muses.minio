package bucketStorage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/silverswords/muses.minio/bucketStorage/driver"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Bucket struct {
	client     Client
	cacher     Cacher
	bucketName string
	ConfigInfo
	mc MinioClient
	uploadTimeOut time.Duration
	maxUploadWorkers int
}

func NewBucket(bucketName, configName, configPath string, uploadTimeOut time.Duration, maxUploadWorkers int) (*Bucket, error) {
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
		uploadTimeOut: uploadTimeOut,
		maxUploadWorkers: maxUploadWorkers,
	}, nil
}

func (b *Bucket) MakeBucket(opts ...driver.OtherMakeBucketOption) error {
	const (
		defaultRegion        = "us-east-1"
		defaultObjectLocking = false
	)

	o := &driver.MakeBucketOptions{
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

func (b *Bucket) SetBucketVersioning(opts ...driver.OtherSetBucketVersioningOption) error {
	const (
		defaultStatus = "Enabled"
	)

	o := &driver.SetBucketVersioningOptions{
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

type ObjectUpload struct {
	Key string
	Object io.Reader
}

type UploadError struct {
	Key string
	Error error
}

type UploadResult struct {
	Key string
}

func (b *Bucket) countNeededWorkers(objectsCount int) int {
	if b.maxUploadWorkers > objectsCount {
		return objectsCount
	}
	 return b.maxUploadWorkers
}

func (b *Bucket) NewTypedWriter(ctx context.Context, objectsCount int, objectChannel chan ObjectUpload, opts ...driver.OtherPutObjectOption) (driver.Writer, error) {
	var e encrypt.ServerSide
	var d = int64(1024)
	var (
		defaultServerSideEncryption = e
		defaultObjectSize = d
	)

	o := &driver.OtherPutObjectOptions{
		defaultServerSideEncryption,
		defaultObjectSize,
	}

	for _, opt := range opts {
		opt(o)
	}

	errorsCh := make(chan *UploadError, objectsCount)

	contextWithTimeout, cancel := context.WithTimeout(ctx, b.uploadTimeOut)
	defer cancel()

	workersCount := b.countNeededWorkers(objectsCount)

	var waitGroup sync.WaitGroup
	waitGroup.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go func() {
			defer waitGroup.Done()
			for {
				select {
				case <-contextWithTimeout.Done():
					return
				default:
				}

				select {
				case upload, ok := <-objectChannel:
					if !ok {
						return
					}
					err := b.client.PutObject(ctx, b.bucketName, upload.Key, upload.Object, o)
					if err != nil {
						errorsCh <- &UploadError{
							upload.Key,
							err,
						}
					}
				default:
				}
			}
		}()
	}

	waitGroup.Wait()
	close(errorsCh)

	return nil, nil
}

func (b *Bucket) SignedURL(ctx context.Context, key string, expires time.Duration, Method string) (string, error) {
	var u *url.URL
	var err error
	switch Method {
	case http.MethodGet:
		u, err = b.client.PresignedGetObject(ctx, b.bucketName, key, expires)
		if err != nil {
			return "", err
		}
		fmt.Println("url", &u)
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

func (b *Bucket) NewRangeReader(ctx context.Context, key string, opts ...driver.OtherGetObjectOption) (driver.Reader, error) {
	var e encrypt.ServerSide
	o := &driver.GetObjectOptions{
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

func (b *Bucket) Delete(ctx context.Context, key string, opts ...driver.OtherRemoveObjectOption) error {
	const (
		defaultGovernanceBypass = false
	)

	o := &driver.RemoveObjectOptions{
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

func (b *Bucket) ListObjects(opts ...driver.OtherListObjectsOption) <-chan minio.ObjectInfo {
	const (
		defaultPrefix = ""
	)

	o := &driver.ListObjectsOptions{
		Prefix: defaultPrefix,
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
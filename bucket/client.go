package bucket

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/spf13/viper"
	"io"
	"os"
)

type Client interface {
	InitClient() error
	PutObject(bucketName string, objectName string, object *os.File, o *OtherPutObjectOptions) error
	GetObject(bucketName string, objectName string, o *OtherGetObjectOptions) ([]byte, error)
	RemoveObject(bucketName string, objectName string, o *OtherRemoveObjectOptions) error
	ListObjects(bucketName string, o *OtherListObjectsOptions) <-chan minio.ObjectInfo
	MakeBucket(bucketName string, o *OtherMakeBucketOptions) error
	CheckBucket(bucketName string) (bool, error)
	ListBuckets() ([]minio.BucketInfo, error)
	RemoveBucket(bucketName string) error
	SetBucketReplication(bucketName string, cfg replication.Config) error
	GetBucketReplication(bucketName string) (replication.Config, error)
	RemoveBucketReplication(bucketName string) error
}

type allConfig struct {
	config map[string]interface{}
}

type minioClient struct {
	mc *minio.Client
}

func (m *minioClient) InitClient() error {
	var b Bucket
	ac, err := b.GetConfig()
	if err != nil {
		return err
	}

	secure := ac.config["secure"]
	endpoint := ac.config["endpoint"]
	accessKeyID := ac.config["access_key_id"]
	secretAccessKey := ac.config["secret_access_key"]

	if endpoint == "" && accessKeyID == "" && secretAccessKey == "" {
		return errors.New("new client failed")
	}

	mc, err := minio.New(endpoint.(string), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID.(string), secretAccessKey.(string), ""),
		Secure: secure.(bool),
	})
	if err != nil {
		return err
	}

	m.mc = mc
	return nil
}

func (m *minioClient) MakeBucket(bucketName string, o *OtherMakeBucketOptions) error {
	err := m.mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{o.Region, o.ObjectLocking})
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) CheckBucket(bucketName string) (bool, error) {
	exists, err := m.mc.BucketExists(context.Background(), bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *minioClient) ListBuckets() ([]minio.BucketInfo, error) {
	bucketInfos, err := m.mc.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func (m *minioClient) RemoveBucket(bucketName string) error {
	err := m.mc.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) SetBucketReplication(bucketName string, cfg replication.Config) error {
	err := m.mc.SetBucketReplication(context.Background(), bucketName, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) GetBucketReplication(bucketName string) (replication.Config, error) {
	cfg, err := m.mc.GetBucketReplication(context.Background(), bucketName)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (m *minioClient) RemoveBucketReplication(bucketName string) error {
	err := m.mc.RemoveBucketReplication(context.Background(), bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) PutObject(bucketName string, objectName string, object *os.File, o *OtherPutObjectOptions) error {
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	_, err = m.mc.PutObject(context.Background(), bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ServerSideEncryption: o.ServerSideEncryption})
	if err != nil {
		return err
	}
	return nil
}

func (m *minioClient) GetObject(bucketName string, objectName string, o *OtherGetObjectOptions) ([]byte, error) {
	minioObject, err := m.mc.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{ServerSideEncryption: o.ServerSideEncryption})
	if err != nil {
		return nil, err
	}

	stat, err := minioObject.Stat()
	buf := make([]byte, stat.Size)
	_, err = io.ReadFull(minioObject, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (m *minioClient) RemoveObject(bucketName string, objectName string, o *OtherRemoveObjectOptions) error {
	err := m.mc.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: o.GovernanceBypass})
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) ListObjects(bucketName string, o *OtherListObjectsOptions) <-chan minio.ObjectInfo {
	objectInfo := m.mc.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Prefix: o.prefix})

	return objectInfo
}

type clientConfigInfo struct {
	configName string
	configPath string
}

func (b *Bucket) GetConfig() (*allConfig, error) {
	var config allConfig
	viper.SetConfigName(b.configName)
	viper.AddConfigPath(b.configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}



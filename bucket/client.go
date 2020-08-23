package bucket

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"io"
	"os"
)

type client interface {
	PutObject(bucketName string, objectName string, object *os.File, opts minio.PutObjectOptions) error
	GetObject(bucketName string, objectName string, opts minio.GetObjectOptions) ([]byte, error)
	RemoveObject(bucketName string, objectName string, opts minio.RemoveObjectOptions) error
	initClient() error
	MakeBucket(bucketName string, opts minio.MakeBucketOptions) error
	CheckBucket(bucketName string) (bool, error)
	ListBuckets() ([]minio.BucketInfo, error)
	RemoveBucket(bucketName string) error
}

type ClientConfig struct {
	Client map[string]interface{}
}

type minioClient struct {
	mc *minio.Client
}

func (m *minioClient) initClient() error {
	var b Bucket
	cc, err := b.getConfig()
	secure := cc.Client["secure"]
	endpoint := cc.Client["endpoint"]
	accessKeyID := cc.Client["access_key_id"]
	secretAccessKey := cc.Client["secret_access_key"]

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

func (m *minioClient) MakeBucket(bucketName string, opts minio.MakeBucketOptions) error {
	err := m.mc.MakeBucket(context.Background(), bucketName, opts)
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

func (m *minioClient) PutObject(bucketName string, objectName string, object *os.File, opts minio.PutObjectOptions) error {
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	_, err = m.mc.PutObject(context.Background(), bucketName, objectName, object, objectStat.Size(), opts)
	if err != nil {
		return err
	}
	return nil
}

func (m *minioClient) GetObject(bucketName string, objectName string, opts minio.GetObjectOptions) ([]byte, error) {
	minioObject, err := m.mc.GetObject(context.Background(), bucketName, objectName, opts)
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

func (m *minioClient) RemoveObject(bucketName string, objectName string, opts minio.RemoveObjectOptions) error {
	err := m.mc.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}

	return nil
}

type clientConfigInfo struct {
	configName string
	configPath string
}

func (b *Bucket) getConfig() (*ClientConfig, error) {
	var config ClientConfig
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



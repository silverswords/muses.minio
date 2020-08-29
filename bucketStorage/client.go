package bucketStorage

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

type Client interface {
	InitClient(configName, configPath string) error
	PutObject(bucketName string, objectName string, object *os.File, o *OtherPutObjectOptions) error
	GetObject(bucketName string, objectName string, o *OtherGetObjectOptions) ([]byte, error)
	RemoveObject(bucketName string, objectName string, o *OtherRemoveObjectOptions) error
	ListObjects(bucketName string, o *OtherListObjectsOptions) <-chan minio.ObjectInfo
	SetObjectLockConfig(bucketName string, mode *minio.RetentionMode, validity *uint, uint *minio.ValidityUnit) error
	GetObjectLockConfig(bucketName string) (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error)
	MakeBucket(bucketName string, o *OtherMakeBucketOptions) error
	CheckBucket(bucketName string) (bool, error)
	ListBuckets() ([]minio.BucketInfo, error)
	RemoveBucket(bucketName string) error
	SetBucketReplication(bucketName string, cfg replication.Config) error
	GetBucketReplication(bucketName string) (replication.Config, error)
	RemoveBucketReplication(bucketName string) error
	SetBucketPolicy(bucketName, policy string) error
	GetBucketPolicy(bucketName string) (string, error)
}

type minioClient struct {
	mc *minio.Client
}

func (m *minioClient) InitClient(configName, configPath string) error {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return err
	}

	log.Println("--------- ac.config ---------", ac.Client)
	secure := ac.Client["secure"]
	endpoint := ac.Client["endpoint"]
	accessKeyID := ac.Client["accessKeyID"]
	secretAccessKey := ac.Client["secretAccessKey"]

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

func (m *minioClient) SetBucketPolicy(bucketName, policy string) error {
	err := m.mc.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) GetBucketPolicy(bucketName string) (string, error) {
	policy, err := m.mc.GetBucketPolicy(context.Background(), bucketName)
	if err != nil {
		return "", err
	}

	return policy, nil
}

func (m *minioClient) SetObjectLockConfig(bucketName string, mode *minio.RetentionMode, validity *uint, uint *minio.ValidityUnit) error {
	err := m.mc.SetObjectLockConfig(context.Background(), bucketName, mode, validity, uint)
	if err != nil {
		return err
	}

	return nil
}

func (m *minioClient) GetObjectLockConfig(bucketName string) (string, *minio.RetentionMode, *uint, *minio.ValidityUnit, error) {
	objectLock, mode, validity, uint, err := m.mc.GetObjectLockConfig(context.Background(), bucketName)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return objectLock, mode, validity, uint, nil
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

type Config struct {
	Client map[string]interface{}
	Cache map[string]interface{}
}

type ConfigInfo struct {
	configName string
	configPath string
}

func GetConfig(configName, configPath string) (*Config, error) {
	var config Config
	log.Println(configName, configPath, "xxxxx")
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	log.Println("aaa")
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}


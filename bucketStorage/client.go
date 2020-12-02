package bucketStorage

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/spf13/viper"
	"io"
	"net/url"
	"time"
)

type Client interface {
	MakeBucket(bucketName string, o *MakeBucketOptions) error
	PutObject(bucketName string, objectName string, reader io.Reader, o *OtherPutObjectOptions) error
	GetObject(bucketName string, objectName string, o *GetObjectOptions) ([]byte, error)
	RemoveObject(bucketName string, objectName string, o *RemoveObjectOptions) error
	ListObjects(bucketName string, o *ListObjectsOptions) <-chan minio.ObjectInfo
}

type MinioClient struct {
	mc *minio.Client
}

func initClient(configName, configPath string) (Client, error) {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return nil, err
	}
	clientType := ac.Client["clientType"]

	if clientType.(string) == "minio" {
		c, err := newMinioClient(configName, configPath)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, nil
}

func newMinioClient(configName, configPath string) (*MinioClient, error) {
	ac, err := GetConfig(configName, configPath)
	if err != nil {
		return nil, err
	}

	secure := ac.Client["secure"]
	endpoint := ac.Client["endpoint"]
	accessKeyID := ac.Client["accessKeyID"]
	secretAccessKey := ac.Client["secretAccessKey"]

	if endpoint == "" && accessKeyID == "" && secretAccessKey == "" {
		return nil, errors.New("new client failed")
	}
	t, err := minio.DefaultTransport(secure.(bool))
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	fmt.Println(endpoint, accessKeyID, secretAccessKey, secure, "parameters")
	mc, err := minio.New(endpoint.(string), &minio.Options{
		Creds:     credentials.NewStaticV4(accessKeyID.(string), secretAccessKey.(string), ""),
		Secure:    secure.(bool),
		Transport: t,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(mc, "mc------------")
	return &MinioClient{
		mc,
	}, nil
}

func (m *MinioClient) MakeBucket(bucketName string, o *MakeBucketOptions) error {
	fmt.Println(bucketName, "bucketName", o.Region, o.ObjectLocking)
	err := m.mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: o.Region, ObjectLocking: o.ObjectLocking})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) CheckBucket(bucketName string) (bool, error) {
	exists, err := m.mc.BucketExists(context.Background(), bucketName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m *MinioClient) ListBuckets() ([]minio.BucketInfo, error) {
	bucketInfos, err := m.mc.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	return bucketInfos, nil
}

func (m *MinioClient) RemoveBucket(bucketName string) error {
	err := m.mc.RemoveBucket(context.Background(), bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) SetBucketVersioning(bucketName string, o *SetBucketVersioningOptions) error {
	err := m.mc.SetBucketVersioning(context.Background(), bucketName, minio.BucketVersioningConfiguration{XMLName: o.XMLName, Status: o.Status, MFADelete: o.MFADelete})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetBucketVersioning(bucketName string) (minio.BucketVersioningConfiguration, error) {
	bucketVersioningConfiguration, err := m.mc.GetBucketVersioning(context.Background(), bucketName)
	if err != nil {
		return bucketVersioningConfiguration, err
	}

	return bucketVersioningConfiguration, nil
}

func (m *MinioClient) SetBucketReplication(bucketName string, cfg replication.Config) error {
	err := m.mc.SetBucketReplication(context.Background(), bucketName, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetBucketReplication(bucketName string) (replication.Config, error) {
	cfg, err := m.mc.GetBucketReplication(context.Background(), bucketName)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (m *MinioClient) RemoveBucketReplication(bucketName string) error {
	err := m.mc.RemoveBucketReplication(context.Background(), bucketName)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) SetBucketPolicy(bucketName, policy string) error {
	err := m.mc.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetBucketPolicy(bucketName string) (string, error) {
	policy, err := m.mc.GetBucketPolicy(context.Background(), bucketName)
	if err != nil {
		return "", err
	}

	return policy, nil
}

func (m *MinioClient) SetObjectLockConfig(bucketName string, mode string, validity *uint, uint string) error {
	mr := (*minio.RetentionMode)(&mode)
	u := (*minio.ValidityUnit)(&uint)
	err := m.mc.SetObjectLockConfig(context.Background(), bucketName, mr, validity, u)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) GetObjectLockConfig(bucketName string) (string, string, *uint, string, error) {
	objectLock, mode, validity, v, err := m.mc.GetObjectLockConfig(context.Background(), bucketName)
	if err != nil {
		return "", "", nil, "", err
	}
	mr := string(*mode)
	u := string(*v)
	return objectLock, mr, validity, u, nil
}

func (m *MinioClient) PutObject(bucketName string, objectName string, reader io.Reader, o *OtherPutObjectOptions) error {
	_, err := m.mc.PutObject(context.Background(), bucketName, objectName, reader, o.ObjectSize, minio.PutObjectOptions{ServerSideEncryption: o.ServerSideEncryption})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioClient) GetObject(bucketName string, objectName string, o *GetObjectOptions) ([]byte, error) {
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

func (m *MinioClient) RemoveObject(bucketName string, objectName string, o *RemoveObjectOptions) error {
	err := m.mc.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{GovernanceBypass: o.GovernanceBypass})
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioClient) PresignedPutObject(bucketName string, objectName string, expires time.Duration) (*url.URL, error) {
	u, err := m.mc.PresignedPutObject(context.Background(), bucketName, objectName, expires)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *MinioClient) PresignedGetObject(bucketName string, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	u, err := m.mc.PresignedGetObject(context.Background(), bucketName, objectName, expires, reqParams)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *MinioClient) ListObjects(bucketName string, o *ListObjectsOptions) <-chan minio.ObjectInfo {
	objectInfo := m.mc.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{Prefix: o.prefix})

	return objectInfo
}

type Config struct {
	Client   map[string]interface{}
	Cache    map[string]interface{}
	Database map[string]interface{}
}

type ConfigInfo struct {
	configName string
	configPath string
}

func GetConfig(configName, configPath string) (*Config, error) {
	var config Config
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
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

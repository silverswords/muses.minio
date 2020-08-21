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
	PutObject(objectName string, object *os.File) error
	GetObject(objectName string) ([]byte, error)
	InitClient(ClientConfig) error
}

type ClientConfig struct {
	Client map[string]interface{}
}

type minioClient struct {
	mc *minio.Client
	bucketName string
}

func (m *minioClient) InitClient(cc *ClientConfig) error {
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

func (m *minioClient) PutObject(objectName string, object *os.File) error {
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	_, err = m.mc.PutObject(context.Background(), m.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *minioClient) GetObject(objectName string) ([]byte, error) {
	minioObject, err := m.mc.GetObject(context.Background(), m.bucketName, objectName, minio.GetObjectOptions{})
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



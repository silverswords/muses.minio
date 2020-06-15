package storage

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	GetMinioClient() *minio.Client
}

type client struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	NewMinioClient  *minio.Client
}

func (c *client) GetMinioClient() *minio.Client {
	return c.NewMinioClient
}

func NewClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool, newMinioClient *minio.Client) Client {
	return &client{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
		NewMinioClient:  newMinioClient,
	}
}

func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) Client {
	newMinioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	minioClient := NewClient(endpoint, accessKeyID, secretAccessKey, useSSL, newMinioClient)
	return minioClient
}

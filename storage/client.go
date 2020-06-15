package storage

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	NewClient()
	GetClient() *minio.Client
}

type client struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Client          *minio.Client
}

func (b *client) NewClient() {
	// Initialize minio client object.
	minioClient, err := minio.New(b.Endpoint, b.AccessKeyID, b.SecretAccessKey, b.UseSSL)
	if err != nil {
		fmt.Println(err)
	}
	b.Client = minioClient
}

func (b *client) GetClient() *minio.Client {
	return b.Client
}

func NewMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) Client {
	return &client{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
	}
}

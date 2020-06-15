package storage

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	New() *minio.Client
}

type client struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func (b *client) New() *minio.Client {
	// Initialize minio client object.
	minioClient, err := minio.New(b.Endpoint, b.AccessKeyID, b.SecretAccessKey, b.UseSSL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return minioClient
}

func NewMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) Client {
	return &client{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		UseSSL:          useSSL,
	}
}

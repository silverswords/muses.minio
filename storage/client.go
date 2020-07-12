package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	getMinioClient() *minio.Client
}

type client struct {
	minioClient *minio.Client
}

type minioClient struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	secure          bool
}

func getMinioClient(m *minioClient) *client {
	newMinioClient, err := minio.New(m.endpoint, m.accessKeyID, m.secretAccessKey, m.secure)
	if err != nil {
		log.Fatalln(err)
	}
	return &client{
		minioClient: newMinioClient,
	}
}

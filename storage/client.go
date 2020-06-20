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

type clientOption func(*client)

var defaultClientOptions = client{
	Endpoint:        "127.0.0.1:9001",
	AccessKeyID:     "minio",
	SecretAccessKey: "minio123",
	UseSSL:          false,
}

func WithEndpoint(e string) clientOption {
	return func(c *client) {
		c.Endpoint = e
	}
}

func WithAccessKeyID(a string) clientOption {
	return func(c *client) {
		c.AccessKeyID = a
	}
}

func WithSecretAccessKey(s string) clientOption {
	return func(c *client) {
		c.SecretAccessKey = s
	}
}

func WithUseSSL(u bool) clientOption {
	return func(c *client) {
		c.UseSSL = u
	}
}

func WithNewMinioClient(n *minio.Client) clientOption {
	return func(c *client) {
		c.NewMinioClient = n
	}
}

func (c *client) GetMinioClient() *minio.Client {
	return c.NewMinioClient
}

func NewClient(options client, newMinioClient *minio.Client) Client {
	return &client{
		Endpoint:        options.Endpoint,
		AccessKeyID:     options.AccessKeyID,
		SecretAccessKey: options.SecretAccessKey,
		UseSSL:          options.UseSSL,
		NewMinioClient:  newMinioClient,
	}
}

func NewMinioClient(opts ...clientOption) Client {
	options := defaultClientOptions
	for _, o := range opts {
		o(&options)
	}

	newMinioClient, err := minio.New(options.Endpoint, options.AccessKeyID, options.SecretAccessKey, options.UseSSL)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	minioClient := NewClient(options, newMinioClient)
	return minioClient
}

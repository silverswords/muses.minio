package storage

import (
	"github.com/minio/minio-go/v6"
)

type Client interface {
	GetMinioClient(EndPoint string) *minio.Client
}

type client struct {
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	NewMinioClient  *minio.Client
}

type clientOption func(*client)

var new Client
var minioClient = new.GetMinioClient("127.0.0.1:9001")

var defaultClientOptions = client{
	AccessKeyID:     "minio",
	SecretAccessKey: "minio123",
	UseSSL:          false,
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

func (c *client) GetMinioClient(EndPoint string) *minio.Client {
	return c.NewMinioClient
}

func NewClient(options client, newMinioClient *minio.Client) Client {
	return &client{
		AccessKeyID:     options.AccessKeyID,
		SecretAccessKey: options.SecretAccessKey,
		UseSSL:          options.UseSSL,
		NewMinioClient:  newMinioClient,
	}
}

func NewMinioClient(EndPoint string, opts ...clientOption) (Client, error) {
	options := defaultClientOptions
	for _, o := range opts {
		o(&options)
	}

	newMinioClient, err := minio.New(EndPoint, options.AccessKeyID, options.SecretAccessKey, options.UseSSL)
	if err != nil {
		return nil, err
	}

	minioClient := NewClient(options, newMinioClient)
	return minioClient, nil
}

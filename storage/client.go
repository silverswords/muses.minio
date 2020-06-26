package storage

import (
	"fmt"
	"log"
	"net/url"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	getMinioClient() *minio.Client
}

type client struct {
	url string
}

type minioClient struct {
	client
	newMinioClient *minio.Client
}

func (m *minioClient) getMinioClient() *minio.Client {
	return m.newMinioClient
}

func newMinioClient(s string) Client {
	useSSl := true

	u, err := url.Parse(s)
	if err != nil {
		log.Fatalln(err)
	}

	p, _ := u.User.Password()

	newMinioClient, err := minio.New(u.Host, u.User.Username(), p, useSSl)
	if err != nil {
		log.Fatalln(err)
	}

	return &minioClient{
		client: client{
			url: s,
		},
		newMinioClient: newMinioClient,
	}
}

func newClient(s string) (Client, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatalln(err)
	}

	if u.Scheme == "minio" {
		return newMinioClient(s), nil
	}
	return nil, fmt.Errorf("Wrong scheme type passed")
}

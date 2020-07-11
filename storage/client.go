package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

// type Client interface {
// 	// getMinioClient() *minio.Client
// }

type client struct {
	minioClient *minio.Client
}

type minioClient struct {
	// url             string `yaml:"url"`
	// endpoint        string `yaml:"endpoint"`
	// accessKeyID     string `yaml:"accessKeyID"`
	// secretAccessKey string `yaml:"secretAccessKey"`
	// secure          bool   `yaml:"secure"`
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	secure          bool
}

func (m *minioClient) getMinioClient() *client {
	// useSSl := true
	// u, err := url.Parse(m.url)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// p, _ := u.User.Password()

	newMinioClient, err := minio.New(m.endpoint, m.accessKeyID, m.secretAccessKey, m.secure)
	if err != nil {
		log.Fatalln(err)
	}
	return &client{
		minioClient: newMinioClient,
	}
}

// func newMinioClient(s string) Client {
// 	return &minioClient{
// 		url: s,
// 	}
// }

// func newClient(s string) (Client, error) {
// 	u, err := url.Parse(s)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	if u.Scheme == "minio" {
// 		return newMinioClient(s), nil
// 	}
// 	return nil, fmt.Errorf("Wrong scheme type passed")
// }

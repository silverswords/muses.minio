package storage

import (
	"fmt"
	"log"
	"net/url"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	setScheme(scheme string)
}

type client struct {
	scheme string
}

type minioClient struct {
	client
	endPoint        string
	accessKeyId     string
	secretAccessKey string
	useSSl          bool
	newMinioClient  *minio.Client
}

var new *minioClient
var m = new.newMinioClient

func (c *client) setScheme(scheme string) {
	c.scheme = scheme
}

func Open(s string) *minioClient {
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
			scheme: u.Scheme,
		},
		endPoint:        u.Host,
		accessKeyId:     u.User.Username(),
		secretAccessKey: p,
		useSSl:          useSSl,
		newMinioClient:  newMinioClient,
	}
}

func GetClient(scheme string, url string) (Client, error) {
	if scheme == "minio" {
		return Open(url), nil
	}
	return nil, fmt.Errorf("Wrong scheme type passed")
}

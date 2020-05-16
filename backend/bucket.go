package backend

import (
	"fmt"
	"io"

	"github.com/minio/minio-go/v6"
)

type minioBackend struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	client          *minio.Client
}

func (b *minioBackend) newMinio() error {
	// Initialize minio client object.
	var err error
	b.client, err = minio.New(b.endpoint, b.accessKeyID, b.secretAccessKey, b.useSSL)
	if err != nil {
		return err
	}
	return nil
}

func (b *minioBackend) checkBucket(bucketName string, location string) error {
	err := b.client.MakeBucket(bucketName, location)
	if err != nil {
		exists, err := b.client.BucketExists(bucketName)
		if err == nil && exists {
			return nil
		}
		return err
	}
	return nil
}

func (b *minioBackend) upload(bucketName string, objectName string, reader io.Reader) (int64, error) {
	n, err := b.client.PutObject(bucketName, objectName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

func (b *minioBackend) read(bucketName string, objectName string) (*minio.Object, error) {
	minioObject, err := b.client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioObject, nil
}

func (b *minioBackend) delete(bucketName string, objectName string) error {
	return b.client.RemoveObject(bucketName, objectName)
}

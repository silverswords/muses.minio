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
	name            string
	bucketName      string
	location        string
	client          *minio.Client
	objectName      string
	reader          io.Reader
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

func (b *minioBackend) checkBucket() error {
	err := b.client.MakeBucket(b.bucketName, b.location)
	if err != nil {
		exists, err := b.client.BucketExists(b.bucketName)
		if err == nil && exists {
			return nil
		}
		return err
	}
	return nil
}

func (b *minioBackend) upload() (int64, error) {
	n, err := b.client.PutObject(b.bucketName, b.objectName, b.reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

func (b *minioBackend) read() (*minio.Object, error) {
	minioObject, err := b.client.GetObject(b.bucketName, b.objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioObject, nil
}

func (b *minioBackend) delete() error {
	return b.client.RemoveObject(b.bucketName, b.objectName)
}

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
	bucketName      string
	location        string
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

func (b *minioBackend) bucketPolicy() error {
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::my-bucketname/*"],"Sid": ""}]}`

	err := b.client.SetBucketPolicy(b.bucketName, policy)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (b *minioBackend) listBuckets() ([]minio.BucketInfo, error) {
	buckets, err := b.client.ListBuckets()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
	return buckets, nil
}

func (b *minioBackend) removeBucket() error {
	return b.client.RemoveBucket(b.bucketName)
}

func (b *minioBackend) upload(objectName string, reader io.Reader) (int64, error) {
	n, err := b.client.PutObject(b.bucketName, objectName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

func (b *minioBackend) read(objectName string) (*minio.Object, error) {
	minioObject, err := b.client.GetObject(b.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioObject, nil
}

func (b *minioBackend) delete(objectName string) error {
	return b.client.RemoveObject(b.bucketName, objectName)
}

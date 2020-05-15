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

func newMinio(name string, info map[string]string) error {
	b := &minioBackend{
		endpoint:        info["endpoint"],
		accessKeyID:     info["access_key_id"],
		secretAccessKey: info["secret_access_key"],
		useSSL:          true,
		name:            name,
		bucketName:      info["bucket_name"],
		location:        info["location"],
		objectName:      info["object_name"],
	}

	if b.endpoint == "" {
		return fmt.Errorf("missing minio param endpoint for %s", name)
	}

	if b.accessKeyID == "" {
		return fmt.Errorf("missing minio param access_key_id for %s", name)
	}

	if b.secretAccessKey == "" {
		return fmt.Errorf("missing minio param secret_access_key for %s", name)
	}

	if b.bucketName == "" {
		return fmt.Errorf("missing minio bucket_name param for %s", name)
	}

	if b.location == "" {
		b.location = "cn-north-1"
	}

	if b.bucketName == "" {
		return fmt.Errorf("missing minio object_name param for %s", name)
	}
	// Initialize minio client object.
	var err error
	b.client, err = minio.New(b.endpoint, b.accessKeyID, b.secretAccessKey, b.useSSL)
	if err != nil {
		return err
	}
	if err = b.checkBucket(); err != nil {
		return err
	}
	if n, err := b.uploadObject(); err != nil {
		fmt.Println(n)
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

func (b *minioBackend) uploadObject() (int64, error) {
	n, err := b.client.PutObject(b.bucketName, b.objectName, b.reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

func (b *minioBackend) Delete() error {
	return b.client.RemoveObject(b.bucketName, b.objectName)
}

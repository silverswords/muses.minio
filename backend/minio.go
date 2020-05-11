package backend

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

func NewMinio(name string, info map[string]string) error {
	endpoint := info["endpoint"]
	accessKeyID := info["access_key_id"]
	secretAccessKey := info["secret_access_key"]

	useSSL := true

	if endpoint == "" {
		return fmt.Errorf("missing minio param endpoint for %s", name)
	}

	if accessKeyID == "" {
		return fmt.Errorf("missing minio param access_key_id for %s", name)
	}

	if secretAccessKey == "" {
		return fmt.Errorf("missing minio param secret_access_key for %s", name)
	}

	b := &minioBackend{
		name:     name,
		bucket:   info["bucket"],
		location: info["location"],
	}

	if b.bucket == "" {
		return fmt.Errorf("missing minio bucket param for %s", name)
	}
	if b.location == "" {
		b.location = "cn-north-1"
	}
	// Initialize minio client object.
	var err error
	b.client, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return err
	}
	if err = b.checkBucket(); err != nil {
		return err
	}
	return nil
}

type minioBackend struct {
	location string
	name     string
	bucket   string
	client   *minio.Client
}

func (b *minioBackend) checkBucket() error {
	err := b.client.MakeBucket(b.bucket, b.location)
	if err != nil {
		exists, err := b.client.BucketExists(b.bucket)
		if err == nil && exists {
			return nil
		}
		return err
	}
	return nil
}

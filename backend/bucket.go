package backend

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

type minioBackend struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	name            string
	bucket          string
	location        string
	client          *minio.Client
}

func newMinio(name string, info map[string]string) error {
	b := &minioBackend{
		endpoint:        info["endpoint"],
		accessKeyID:     info["access_key_id"],
		secretAccessKey: info["secret_access_key"],
		useSSL:          true,
		name:            name,
		bucket:          info["bucket"],
		location:        info["location"],
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

	if b.bucket == "" {
		return fmt.Errorf("missing minio bucket param for %s", name)
	}

	if b.location == "" {
		b.location = "cn-north-1"
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
	return nil
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

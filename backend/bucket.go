package backend

import (
	"fmt"
	"net/url"
	"time"

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
	expiry          time.Duration
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
		expiry:          time.Second * 24 * 60 * 60, // Generates a url which expires in a day.
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
	if url, err := b.createObject(); err != nil {
		fmt.Println(url)
		return err
	}
	return nil
}

func (b *minioBackend) createObject() (*url.URL, error) {
	presignedURL, err := b.client.PresignedPutObject(b.bucketName, b.objectName, b.expiry)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return presignedURL, nil
}

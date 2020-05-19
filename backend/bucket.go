package backend

import (
	"fmt"
	"io"
	"log"

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
	mode            *minio.RetentionMode
	validity        *uint
	unit            *minio.ValidityUnit
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
			log.Printf("We already own %s\n", b.bucketName)
		}
	}
	return err
}

func (b *minioBackend) removeBucket() error {
	return b.client.RemoveBucket(b.bucketName)
}

func (b *minioBackend) bucketPolicy() error {
	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::my-bucketname/*"],"Sid": ""}]}`

	err := b.client.SetBucketPolicy(b.bucketName, policy)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (b *minioBackend) objectLock() error {
	// mode := Governance
	// validity := uint(30)
	// unit := Days

	err := b.client.SetObjectLockConfig(b.bucketName, b.mode, b.validity, b.unit)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (b *minioBackend) bucketEncryption() error {
	config := minio.ServerSideEncryptionConfiguration{Rules: []minio.Rule{
		{
			Apply: minio.ApplyServerSideEncryptionByDefault{
				SSEAlgorithm: "AES256",
			},
		},
	}}
	// Set default encryption configuration on a bucket
	err := b.client.SetBucketEncryption(b.bucketName, config)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (b *minioBackend) uploadFile(filePath string, reader io.Reader) (int64, error) {
	n, err := b.client.PutObject(b.bucketName, filePath, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

func (b *minioBackend) readFile(filePath string) (*minio.Object, error) {
	minioObject, err := b.client.GetObject(b.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioObject, nil
}

func (b *minioBackend) deleteFile(filePath string) error {
	return b.client.RemoveObject(b.bucketName, filePath)
}

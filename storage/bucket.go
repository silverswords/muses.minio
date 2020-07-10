package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName string
	location   string
	client
	bucketObjectCache
}

func NewBucket(bucketName, location string) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		location:   location,
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
	}
}

func (b *Bucket) MakeBucket() error {
	err := b.client.minioClient.MakeBucket(b.bucketName, b.location)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (b *Bucket) CheckBucket() (bool, error) {
	exists, err := b.client.minioClient.BucketExists(b.bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName  string
	location    string
	clients     []*minio.Client
	minioClient *minio.Client
	bucketObjectCache
	Weight
}

func NewBucket(bucketName, location string) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		location:   location,
		clients:    getMinioClients(),
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
	}
}

func (b *Bucket) MakeBucket() error {
	for _, v := range b.clients {
		err := v.MakeBucket(b.bucketName, b.location)
		if err != nil {
			log.Fatalln(err)
			return err
		}
	}
	return nil
}

func (b *Bucket) CheckBucket() (bool, error) {
	exists, err := b.minioClient.BucketExists(b.bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

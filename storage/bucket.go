package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName            string
	location              string
	minioClientWithWeight map[string]strategyClient
	bucketObjectCache
	strategyClients []*strategyClient
	strategy        string
}

func NewBucket(bucketName, location string, strategy string) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		location:   location,
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
		strategyClients: getStrategyClients(),
		strategy:        strategy,
	}
}

func (b *Bucket) MakeBucket() error {
	for _, v := range b.strategyClients {
		err := v.client.MakeBucket(b.bucketName, b.location)
		if err != nil {
			log.Fatalln(err)
			return err
		}
	}
	return nil
}

func (b *Bucket) CheckBucket() (bool, error) {
	var exists bool
	for _, v := range b.strategyClients {
		err := v.client.MakeBucket(b.bucketName, b.location)
		exists, err = v.client.BucketExists(b.bucketName)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return exists, nil
}

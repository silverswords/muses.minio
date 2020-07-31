package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName string
	location   string
	// minioClientWithWeight map[string]strategyClient
	bucketObjectCache
	strategy string
}

func NewBucket(bucketName, location string, strategy string) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		location:   location,
		// bucketObjectCache: bucketObjectCache{
		// 	items: make(map[string]*minio.Object),
		// },
		// strategyClients: getStrategyClients(),
		strategy: strategy,
	}
}

func (b *Bucket) MakeBucket() error {
	for _, v := range getStrategyClients() {
		err := v.client.MakeBucket(b.bucketName, b.location)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (b *Bucket) MakeBucket() error {
// 	for _, v := range b.strategyClients {
// 		err := v.client.MakeBucket(b.bucketName, b.location)
// 		if err != nil {
// 			log.Fatalln(err)
// 			return err
// 		}
// 	}
// 	return nil
// }

func (b *Bucket) CheckBucket(bucketName string) (bool, error) {
	var exists bool
	var err error
	for _, v := range getStrategyClients() {
		exists, err = v.client.BucketExists(bucketName)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return exists, nil
}

func (b *Bucket) ListedBucket() []minio.BucketInfo {
	minioClient := getStrategyClients()[0].client
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		log.Fatalln(err)
	}

	return buckets
}

func (b *Bucket) RemoveBucket(bucketName string) error {
	for _, v := range getStrategyClients() {
		err := v.client.RemoveBucket(bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type StrategyClient interface {
	Save(minio.Object) error
}

// 你这个 bucket 不是它定义的 bucket，你是它 bucket 更上一层的抽象
type Bucket struct {
	bucketName string
	location   string
	// minioClients []*minio.Client
	minioClientWithWeight map[string]strategyClient
	bucketObjectCache
	strategyClients []*strategyClient
	strategy        string
}

type strategyClient struct {
	client *minio.Client
	weight float64
}

func newStrategyClient(client *minio.Client, weight float64) *strategyClient {
	return &strategyClient{
		client: client,
		weight: weight,
	}
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

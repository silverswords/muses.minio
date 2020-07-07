package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName        string `yaml:"bucketName"`
	location          string
	client            `yaml:"client"`
	bucketObjectCache `yaml:"bucketObjectCache"`
}

// func newBucket(s string, bucketName string) *Bucket {
// 	return &Bucket{
// 		bucketName: bucketName,
// 		client: client{
// 			minioClient: minioClient{
// 				url: s,
// 			},
// 		},
// 		bucketObjectCache: bucketObjectCache{
// 			items: make(map[string]*minio.Object),
// 		},
// 	}
// }

func newBucket(s string, bucketName string) *Bucket {
	return &Bucket{
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
	}
}

func (b *Bucket) makeBucket() error {
	err := b.client.getMinioClient().MakeBucket(b.bucketName, b.location)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (b *Bucket) checkBucket() (bool, error) {
	exists, err := b.client.getMinioClient().BucketExists(b.bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

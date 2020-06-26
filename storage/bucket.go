package storage

import (
	"log"
	"net/url"

	"github.com/minio/minio-go/v6"
)

type Bucket struct {
	bucketName string
	minioClient
	bucketObjectCache
}

func newBucket(s string, bucketName string) *Bucket {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatalln(err)
	}
	nc, _ := newClient(u.Scheme)
	mc := nc.getMinioClient()

	return &Bucket{
		bucketName: bucketName,
		minioClient: minioClient{
			client: client{
				url: s,
			},
			newMinioClient: mc,
		},
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
	}
}

func (b *Bucket) makeBucket() error {
	location := "cn-north-1"
	err := b.newMinioClient.MakeBucket(b.bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (b *Bucket) checkBucket() (bool, error) {
	exists, err := b.newMinioClient.BucketExists(b.bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

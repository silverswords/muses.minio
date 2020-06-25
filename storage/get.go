package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) Get(bucketName string, objectName string) (*minio.Object, error) {
	minioObject := objectCache.Get(bucketName, objectName)
	if minioObject == nil {
		minioObject, err := m.GetObject(bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}

		err = objectCache.set(bucketName, objectName)
		if err != nil {
			log.Fatalln(err)
		}

		return minioObject, nil
	}

	return minioObject, nil
}

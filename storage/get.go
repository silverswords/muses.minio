package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) Get(objectName string) (*minio.Object, error) {
	minioObject := b.cacheGet(objectName)
	if minioObject == nil {
		minioObject, err := b.newMinioClient.GetObject(b.bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}

		err = b.cacheSave(objectName)
		if err != nil {
			log.Fatalln(err)
		}

		return minioObject, nil
	}

	return minioObject, nil
}

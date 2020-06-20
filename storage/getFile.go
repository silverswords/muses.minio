package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

func GetFile(bucketName string, objectName string) (*minio.Object, error) {
	new, err := NewMinioClient()
	if err != nil {
		return nil, err
	}
	minioClient := new.GetMinioClient()

	var b *bucketObjectCache
	minioObject := b.Get(bucketName, objectName)
	if minioObject == nil {
		minioObject, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}

		err = b.Set(bucketName, objectName)
		if err != nil {
			log.Fatalln(err)
		}

		return minioObject, nil
	}

	return minioObject, nil
}

package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) GetObject(bucketName string, objectName string) (*minio.Object, error) {
	var err error
	minioObject := b.cacheGet(objectName)
	if minioObject == nil {
		for i := 0; i < len(getStrategyClients()); i++ {
			minioObject, err = getStrategyClients()[i].client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			if minioObject != nil {
				break
			}
		}

		err := b.cacheSave(bucketName, objectName)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return minioObject, nil
}

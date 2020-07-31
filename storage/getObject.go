package storage

import (
	"github.com/minio/minio-go/v6"
)

func (b *Bucket) GetObject(objectName string) (*minio.Object, error) {
	var err error
	var minioObject *minio.Object
	// minioObject := b.cacheGet(objectName)
	// if minioObject == nil {
	// 	for i := 0; i < len(getStrategyClients()); i++ {
	// 		minioObject, err = getStrategyClients()[i].client.GetObject(b.bucketName, objectName, minio.GetObjectOptions{})
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if minioObject != nil {
	// 			break
	// 		}
	// 	}

	// 	b.cacheSave(objectName, minioObject)
	// }
	for i := 0; i < len(getStrategyClients()); i++ {
		minioObject, err = getStrategyClients()[i].client.GetObject(b.bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		if minioObject != nil {
			break
		}
	}

	return minioObject, nil
}

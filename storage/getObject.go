package storage

import (
	"fmt"
	"github.com/minio/minio-go/v6"
)

func (b *Bucket) GetObject(objectName string) (*minio.Object, error) {
	var err error
	var minioObject *minio.Object
	var buf []byte
	buf = b.getCacheObject(objectName)
	fmt.Println(buf)

	for i := 0; i < len(getStrategyClients()); i++ {
		fmt.Println(b.bucketName, objectName, "bucketName objectName")
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

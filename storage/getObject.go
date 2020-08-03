package storage

import (
	"github.com/minio/minio-go/v6"
	"io"
	"log"
)

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	var err error
	var minioObject *minio.Object
	var buf []byte
	buf = b.getCacheObject(objectName)

	if buf == nil {
		var object []byte
		for i := 0; i < len(b.getStrategyClients()); i++ {
			minioObject, err = b.getStrategyClients()[i].client.GetObject(b.bucketName, objectName, minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			if minioObject != nil {
				break
			}
		}
		_, err = io.ReadFull(minioObject, object)
		if err != nil {
			log.Fatalln(err)
		}

		buf = object
	}


	return buf, nil
}

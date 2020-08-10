package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	var minioObject *minio.Object
	var buf []byte
	if b.cache {
		buf, err := b.getCacheObject(objectName)
		if err == nil && buf != nil {
			return buf, nil
		}
	}

	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	if buf == nil {
		for _, v := range clients {
			minioObject, err = v.client.GetObject(context.Background(), b.bucketName, objectName, minio.GetObjectOptions{})
			if err != nil {
				return nil, err
			}
			if minioObject != nil {
				break
			}
		}

		stat, err := minioObject.Stat()
		buf = make([]byte, stat.Size)
		_, err = io.ReadFull(minioObject, buf)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

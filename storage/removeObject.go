package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func (b *Bucket) RemoveObject(objectName string) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.RemoveObject(context.Background(), b.bucketName, objectName, minio.RemoveObjectOptions{true, ""})
		if err != nil {
			return err
		}
	}

	if b.cache {
		err = b.deleteCacheObject(objectName)
		if err != nil {
			return err
		}
	}

	return nil
}

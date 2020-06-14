package storage

import (
	"fmt"

	"github.com/minio/minio-go"
)

func (b *minioBackend) Get(filePath string) (*minio.Object, error) {
	minioObject, err := b.client.GetObject(b.bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return minioObject, nil
}

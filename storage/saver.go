package storage

import (
	"fmt"
	"io"

	"github.com/minio/minio-go/pkg/encrypt"
	"github.com/minio/minio-go/v6"
)

type Saver interface {
	Save()
}

func putOptions(encrypted bool, contentType string, mode *minio.RetentionMode) minio.PutObjectOptions {
	options := minio.PutObjectOptions{}
	if encrypted {
		options.ServerSideEncryption = encrypt.NewSSE()
	}
	options.ContentType = contentType
	options.Mode = mode

	return options
}

func (b *minioBackend) Save(filePath string, reader io.Reader) (int64, error) {
	contentType := "binary/octet-stream"

	options := putOptions(b.encrypt, contentType, b.mode)
	n, err := b.client.PutObject(b.bucketName, filePath, reader, -1, options)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return n, nil
}

package storage

import (
	"io"

	"github.com/minio/minio-go"
)

type backend interface {
	Save(filePath string, reader io.Reader) (int64, error)
	Get(filePath string) (*minio.Object, error)
	Delete(filePath string) error
}

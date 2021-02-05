package driver

import (
	"context"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"time"
)

type Bucket interface {
	NewRangeReader(ctx context.Context, key string, opts ...bucketStorage.OtherGetObjectOption) (Reader, error)
	NewTypedWriter(ctx context.Context, key string, opts ...bucketStorage.OtherPutObjectOption) (Writer, error)
	Delete(ctx context.Context, key string, opts ...bucketStorage.OtherRemoveObjectOption) error
	SignedURL(ctx context.Context, key string, expires time.Duration, Method string) (string, error)
	Close() error
}

type Reader interface {
	io.ReadCloser
}

type Writer interface {
	io.WriteCloser
}


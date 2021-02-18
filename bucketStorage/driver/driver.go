package driver

import (
	"context"
	"io"
	"time"
)

type Bucket interface {
	NewRangeReader(ctx context.Context, key string, opts ...OtherGetObjectOption) (Reader, error)
	NewTypedWriter(ctx context.Context, key string, reader io.Reader, opts ...OtherPutObjectOption) (Writer, error)
	Delete(ctx context.Context, key string, opts ...OtherRemoveObjectOption) error
	SignedURL(ctx context.Context, key string, expires time.Duration, Method string) (string, error)
	Close() error
}

type Reader interface {
	io.ReadCloser
}

type Writer interface {
	io.WriteCloser
}


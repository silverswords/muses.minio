package driver

import (
	"context"
	"io"
	"time"
)

type Bucket interface {
	NewRangeReader(ctx context.Context, key string, offset, length int64, opts *ReaderOptions) (Reader, error)
	NewTypedWriter(ctx context.Context, key, contentType string, opts *WriterOptions) (Writer, error)
	Delete(ctx context.Context, key string) error
	SignedURL(ctx context.Context, key string, opts *SignedURLOptions) (string, error)
	ListPaged(ctx context.Context, opts *ListOptions) (*ListPage, error)
	Close() error
}

type Reader interface {
	io.ReadCloser
}

type Writer interface {
	io.WriteCloser
}

type ReaderOptions struct {
	BeforeRead func(asFunc func(interface{}) bool) bool
}

type WriterOptions struct {
	BufferSize int
	// CacheControl specifies caching attributes that services may use
	// when serving the blob.
	CacheControl string
	// ContentMD5 is used as a message integrity check.
	ContentMD5 []byte
	Metadata map[string]string
	BeforeWrite func(asFunc func(interface{}) bool) bool
}

type SignedURLOptions struct {
	Expiry time.Duration
	Method string
}

type ListOptions struct {
	Prefix string
	Delimiter string
	BeforeList func(asFunc func(interface{}) bool) bool
}

type ListObject struct {
	Key string
	Size int64
	MD5 []byte
}

type ListPage struct {
	Object []*ListObject
	NextPageToken []byte
}

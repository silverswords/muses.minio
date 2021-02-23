package main

import (
	"context"
	"github.com/silverswords/muses.minio/bucketStorage"
	"github.com/silverswords/muses.minio/bucketStorage/driver"
	"github.com/silverswords/muses.minio/bucketStorage/middleware"
	"io"
)

type writerOptions struct {
	key string
	file io.Reader
	opts driver.OtherPutObjectOption
}

func NewWriter(ctx context.Context, mn interface{}) (interface{}, error) {
	var bucket *bucketStorage.Bucket
	if m, ok := mn.(*writerOptions); ok {
		_, err := bucket.NewTypedWriter(ctx, m.key, m.file)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func main() {
	judge := make(map[string]bool)
	judge["limit"] = true

	cond := &middleware.Condition{
		ObjectSize: 20,
		AllSize: 40,
		Threshold: 100,
	}

	middleware.Chain(judge, middleware.Limit(cond))(NewWriter)
}

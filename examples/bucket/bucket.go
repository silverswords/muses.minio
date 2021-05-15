package main

import (
	"context"
	"io"

	"github.com/silverswords/muses.minio/bucketStorage"
	"github.com/silverswords/muses.minio/bucketStorage/driver"
	"github.com/silverswords/muses.minio/bucketStorage/middleware"
)

type writerOptions struct {
	key         string
	objectCount int
	file        io.Reader
	opts        driver.OtherPutObjectOption
}

func NewWriter(ctx context.Context, mn interface{}) (interface{}, error) {
	var bucket *bucketStorage.Bucket
	var uploadChan chan *bucketStorage.ObjectUpload
	if m, ok := mn.(*writerOptions); ok {
		uploadChan <- &bucketStorage.ObjectUpload{Key: m.key, Object: m.file}
		_, err := bucket.NewTypedWriter(ctx, m.objectCount, uploadChan)
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
		AllSize:    40,
		Threshold:  100,
	}

	middleware.Chain(judge, middleware.Limit(cond))(NewWriter)
}

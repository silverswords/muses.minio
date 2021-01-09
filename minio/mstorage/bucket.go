package mstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/silverswords/muses.minio/minio"
	"github.com/silverswords/muses.minio/minio/driver"
	"io"
	"net/http"
	"net/url"
)

type clientConfig struct {
	endpoint string
	accessKeyID string
	secretAccessKey string
	useSSL bool
}

func openBucket(ctx context.Context, conf clientConfig, bucketName string) (*bucket, error) {
	if bucketName == "" {
		return nil, errors.New("mstorage.OpenBucket: bucketName is required")
	}

	minioClient, err := minio.New(conf.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.accessKeyID, conf.secretAccessKey, ""),
		Secure: conf.useSSL,
	})
	if err != nil {
		return nil, err
	}
	
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return nil, errors.New("mstorage.OpenBucket: bucket creation failed")
	}
	return &bucket{
		name:          bucketName,
		client:        minioClient,
	}, nil
}

func OpenBucket(ctx context.Context, conf clientConfig, bucketName string) (*mstorage.Bucket, error) {
	drv, err := openBucket(ctx, conf, bucketName)
	if err != nil {
		return nil, err
	}
	return mstorage.NewBucket(drv), nil
}

type bucket struct {
	name string
	client *minio.Client
}

type writer struct {
	w *io.PipeWriter // created when the first byte is written

	ctx      context.Context
	uploader *minio.UploadInfo
	donec    chan struct{} // closed when done writing
	// The following fields will be written before donec closes:
	err error
}

// Write appends p to w. User must call Close to close the w after done writing.
func (w *writer) Write(p []byte) (int, error) {
	// Avoid opening the pipe for a zero-length write;
	// the concrete can do these for empty blobs.
	if len(p) == 0 {
		return 0, nil
	}
	if w.w == nil {
		// We'll write into pw and use pr as an io.Reader for the
		// Upload call to S3.
		pr, pw := io.Pipe()
		w.w = pw
		if err := w.open(pr); err != nil {
			return 0, err
		}
	}
	select {
	case <-w.donec:
		return 0, w.err
	default:
	}
	return w.w.Write(p)
}

// pr may be nil if we're Closing and no data was written.
func (w *writer) open(pr *io.PipeReader) error {
	return nil
}

// Close completes the writer and closes it. Any error occurring during write
// will be returned. If a writer is closed before any Write is called, Close
// will create an empty file at the given key.
func (w *writer) Close() error {
	if w.w == nil {
		// We never got any bytes written. We'll write an http.NoBody.
		w.open(nil)
	} else if err := w.w.Close(); err != nil {
		return err
	}
	<-w.donec
	return w.err
}

func (b *bucket) Close() error {
	return nil
}

func (b *bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opt *driver.ReaderOptions) (driver.Reader, error) {
	opts := minio.GetObjectOptions{}
	object, err := b.client.GetObject(ctx, b.name, key, opts)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (b *bucket) NewTypedWriter(ctx context.Context, key string, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	r := bytes.NewReader(opts.ContentMD5)
	info, err := b.client.PutObject(ctx, b.name, key, r, int64(opts.BufferSize), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	return &writer{
		ctx:      ctx,
		uploader: &info,
		donec:    make(chan struct{}),
	}, nil
}

func (b *bucket) Delete(ctx context.Context, key string) error {
	opts := minio.RemoveObjectOptions {
		GovernanceBypass: true,
	}
	err := b.client.RemoveObject(ctx, b.name, key, opts)
	return err
}

func (b *bucket) SignedURL(_ context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	var u *url.URL
	var err error
	switch opts.Method {
	case http.MethodGet:
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename=\"file\"")
		u, err = b.client.PresignedGetObject(context.Background(), b.name, key, opts.Expiry, reqParams)
		if err != nil {
			return "", err
		}
	case http.MethodPut:
		u, err = b.client.PresignedPutObject(context.Background(), b.name, key, opts.Expiry)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unsupported Method #{opts.Method}")
	}
	return u.Path, nil
}

func (b *bucket) listObjects(ctx context.Context, opts *driver.ListOptions) <- chan minio.ObjectInfo {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := b.client.ListObjects(ctx, b.name, minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
		}
		fmt.Println(object)
	}
	return objectCh
}

func (b *bucket) ListPaged(ctx context.Context, opts *driver.ListOptions) (*driver.ListPage, error) {
	//var object []*driver.ListObject
	_ = b.listObjects(ctx, opts)
	page := driver.ListPage{}
	return &page, nil
}

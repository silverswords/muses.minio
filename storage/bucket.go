package storage

import "log"

type bucket struct {
	Location string
}

type Option func(*bucket)

var defaultOptions = bucket{
	Location: "cn-north-1",
}

func WithLocation(l string) Option {
	return func(b *bucket) {
		b.Location = l
	}
}

func checkBucket(bucketName string) (bool, error) {
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

func newBucket(bucketName string, opts ...Option) error {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}

	err := minioClient.MakeBucket(bucketName, options.Location)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return err
}

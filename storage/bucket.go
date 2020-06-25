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

func (m *minioClient) CheckBucket(bucketName string) (bool, error) {
	minioClient := m.newMinioClient
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

func (m *minioClient) NewBucket(bucketName string, opts ...Option) error {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}

	minioClient := m.newMinioClient
	err := minioClient.MakeBucket(bucketName, options.Location)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return err
}

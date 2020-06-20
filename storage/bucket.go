package storage

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

func checkBucket(bucketName string, opts ...Option) error {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}

	new, err := NewMinioClient()
	if err != nil {
		return err
	}
	minioClient := new.GetMinioClient()

	exists, err := minioClient.BucketExists(bucketName)
	if exists == false && err == nil {
		err = minioClient.MakeBucket(bucketName, options.Location)
		if err != nil {
			return err
		}
	}

	return err
}

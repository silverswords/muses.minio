package bucket

import "os"

type Bucket struct {
	client  client
	cache   cacher
	clientConfigInfo
	minioClient
	cacheObject
}

type bucketObject struct {
	objectName string
	object *os.File
}

func (b *bucketObject) PutObject() error {
	return nil
}

func (b *bucketObject) GetObject() ([]byte, error) {
	var buf []byte
	return buf, nil
}

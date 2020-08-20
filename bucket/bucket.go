package bucket

import "os"

type Bucket struct {
	client  client
	cache   cache
	bucketName string
}

func (b *Bucket) PutObject(objectName string, object *os.File) error {
	return nil
}

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	var buf []byte
	return buf, nil
}

package bucket

import "os"

type Bucket struct {
	client  client
	cache   cacher
	clientConfigInfo
	minioClient
	cacheObject
}

func (b *Bucket) PutObject(objectName string, object *os.File) error {
	err := b.client.PutObject(objectName, object)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bucket) GetObject(objectName string) ([]byte, error) {
	var buf []byte
	buf, err := b.client.GetObject(objectName)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

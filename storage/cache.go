package storage

import (
	"sync"

	"github.com/minio/minio-go/v6"
)

type bucketObjectCache struct {
	// mutex is used for handling the concurrent
	// read/write requests for cache.
	sync.RWMutex

	items map[string]*minio.Object
}

func newBucketObjectCache() *bucketObjectCache {
	return &bucketObjectCache{
		items: make(map[string]*minio.Object),
	}
}

func (b *bucketObjectCache) Get(bucketName string, objectName string) (minioObject *minio.Object) {
	b.RLock()
	defer b.RUnlock()

	filePath := bucketName + "/" + objectName
	minioObject = b.items[filePath]

	return minioObject
}

func (b *bucketObjectCache) Set(bucketName string, objectName string) error {
	b.Lock()
	defer b.Unlock()

	new, err := NewMinioClient()
	if err != nil {
		return err
	}
	minioClient := new.GetMinioClient()

	filePath := bucketName + "/" + objectName
	minioObject, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	b.items[filePath] = minioObject

	return err
}

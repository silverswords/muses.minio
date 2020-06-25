package storage

import (
	"sync"

	"github.com/minio/minio-go/v6"
)

var objectCache *bucketObjectCache

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

func (bc *bucketObjectCache) Get(bucketName string, objectName string) *minio.Object {
	bc.RLock()
	defer bc.RUnlock()

	filePath := bucketName + "/" + objectName
	minioObject := bc.items[filePath]

	return minioObject
}

func (bc *bucketObjectCache) set(bucketName string, objectName string) error {
	bc.Lock()
	defer bc.Unlock()

	bk := newBucket(bucketName)
	filePath := bucketName + "/" + objectName
	minioObject, err := bk.Get(bucketName, objectName)
	bc.items[filePath] = minioObject

	return err
}

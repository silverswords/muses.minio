package bucket

type cache interface {
	setCacheObject(minioObject []byte, objectName string) error
	getCacheObject(objectName string) ([]byte, error)
}

type objectCache struct {
}

func (o *objectCache) setCacheObject(minioObject []byte, objectName string) error {
	return nil
}

func (o *objectCache) getCacheObject(objectName string) ([]byte, error) {
	var buf []byte
	return buf, nil
}
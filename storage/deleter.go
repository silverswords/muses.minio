package storage

type Deleter interface {
	Delete(filePath string) error
}

func (b *minioBackend) Delete(filePath string) error {
	return b.client.RemoveObject(b.bucketName, filePath)
}

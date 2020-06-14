package storage

func (b *minioBackend) Delete(filePath string) error {
	return b.client.RemoveObject(b.bucketName, filePath)
}

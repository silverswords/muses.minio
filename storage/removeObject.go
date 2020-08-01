package storage

func (b *Bucket) RemoveObject(objectName string) error {
	for _, v := range getStrategyClients() {
		err := v.client.RemoveObject(b.bucketName, objectName)
		if err != nil {
			return err
		}
	}

	b.deleteCacheObject(objectName)

	return nil
}

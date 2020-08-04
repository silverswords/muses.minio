package storage

func (b *Bucket) RemoveObject(objectName string) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err := v.client.RemoveObject(b.bucketName, objectName)
		if err != nil {
			return err
		}
	}

	err = b.deleteCacheObject(objectName)
	if err != nil {
		return err
	}

	return nil
}

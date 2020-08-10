package storage

import (
	"context"
	"github.com/minio/minio-go/v7/pkg/sse"
)

func (b *Bucket) setBucketEncryption() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SetBucketEncryption(context.Background(), b.bucketName, sse.NewConfigurationSSES3())
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketEncryption() (*sse.Configuration, error) {
	var encryptionConfig *sse.Configuration
	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	for _, v := range clients {
		encryptionConfig, err = v.client.GetBucketEncryption(context.Background(), b.bucketName)
		if err != nil {
			return nil, err
		}

		if encryptionConfig != nil {
			break
		}
	}

	return encryptionConfig, nil
}

func (b *Bucket) removeBucketEncryption() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.RemoveBucketEncryption(context.Background(), b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

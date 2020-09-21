package storage

import (
	"context"
	"github.com/minio/minio-go/v7/pkg/notification"
)

func (b *Bucket) setBucketNotification(config notification.Configuration) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SetBucketNotification(context.Background(), b.bucketName, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) listenBucketNotification(ctx context.Context, prefix, suffix string, events []string) (<-chan notification.Info, error) {
	var notificationInfo <-chan notification.Info
	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	for _, v := range clients {
		notificationInfo = v.client.ListenBucketNotification(ctx, b.bucketName, prefix, suffix, events)
	}

	return notificationInfo, nil
}

func (b *Bucket) getBucketNotification() (notification.Configuration, error) {
	var config notification.Configuration
	clients, err := b.getStrategyClients()
	if err != nil {
		return config, err
	}

	for _, v := range clients {
		config, err = v.client.GetBucketNotification(context.Background(), b.bucketName)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func (b *Bucket) removeBucketAllNotification() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.RemoveAllBucketNotification(context.Background(), b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

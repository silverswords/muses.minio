package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func (b *Bucket) enableBucketVersioning() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.EnableVersioning(context.Background(), b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketVersioning() (minio.BucketVersioningConfiguration, error) {
	clients, err := b.getStrategyClients()

	var vc, versionConfig minio.BucketVersioningConfiguration
	if err != nil {
		return versionConfig, err
	}

	for _, v := range clients {
		vc, err = v.client.GetBucketVersioning(context.Background(), b.bucketName)
		if err != nil {
			return versionConfig, err
		}
	}

	return vc, nil
}

func (b *Bucket) suspendBucketVersioning() error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SuspendVersioning(context.Background(), b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}

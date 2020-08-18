package storage

import (
	"github.com/minio/minio-go/v7/pkg/replication"
	"context"
)

func (b *Bucket) setBucketReplication(ctx context.Context, cfg replication.Config) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SetBucketReplication(ctx, b.bucketName, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketReplication(ctx context.Context) (replication.Config, error) {
	var cfg replication.Config
	clients, err := b.getStrategyClients()
	if err != nil {
		return cfg, err
	}

	for _, v := range clients {
		cfg, err = v.client.GetBucketReplication(ctx, b.bucketName)
		if err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

func (b *Bucket) removeBucketReplication(ctx context.Context) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.RemoveBucketReplication(ctx, b.bucketName)
		if err != nil {
			return err
		}
	}

	return nil
}




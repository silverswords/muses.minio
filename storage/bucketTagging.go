package storage

import (
	"context"
	"github.com/minio/minio-go/v7/pkg/tags"
)

func (b *Bucket) setBucketTagging(ctx context.Context, tags *tags.Tags) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SetBucketTagging(ctx, b.bucketName, tags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketTagging(ctx context.Context) (*tags.Tags, error) {
	var tag *tags.Tags
	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	for _, v := range clients {
		tag, err = v.client.GetBucketTagging(ctx, b.bucketName)
		if err != nil {
			return nil, err
		}

		if tag != nil {
			break
		}
	}

	return tag, nil
}

func (b *Bucket) removeBucketTagging(ctx context.Context) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.RemoveBucketTagging(ctx, b.bucketName)
		if err != nil {
			return err
		}
	}
	return nil
}

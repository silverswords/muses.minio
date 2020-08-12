package storage

import (
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"context"
)

// Set lifecycle on a bucket
func (b *Bucket) setBucketLifecycle(days lifecycle.ExpirationDays) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "rule1",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: days,
			},
		},
	}
	for _, v := range clients {
		err = v.client.SetBucketLifecycle(context.Background(), b.bucketName, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketLifecycle() (*lifecycle.Configuration ,error) {
	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	var lc *lifecycle.Configuration
	for _, v := range clients {
		lc, err = v.client.GetBucketLifecycle(context.Background(), b.bucketName)
		if err != nil {
			return nil, err
		}

		if lc != nil {
			break
		}
	}

	return lc, nil
}

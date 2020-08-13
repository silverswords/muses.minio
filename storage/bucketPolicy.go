package storage

import (
	"context"
)

func (b *Bucket) setBucketPolicy(policy string) error {
	clients, err := b.getStrategyClients()
	if err != nil {
		return err
	}

	for _, v := range clients {
		err = v.client.SetBucketPolicy(context.Background(), b.bucketName, policy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bucket) getBucketPolicy() (string, error) {
	clients, err := b.getStrategyClients()
	if err != nil {
		return "", err
	}

	var policy string
	for _, v := range clients {
		policy, err = v.client.GetBucketPolicy(context.Background(), b.bucketName)
		if err != nil {
			return "", err
		}

		if policy != "" {
			break
		}
	}

	return policy, nil
}

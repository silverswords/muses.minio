package storage

import (
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type strategyClient struct {
	client *minio.Client
	weight float64
}

func newStrategyClient(client *minio.Client, weight float64) *strategyClient {
	return &strategyClient{
		client: client,
		weight: weight,
	}
}

func (b *Bucket) getStrategyClients() ([]*strategyClient, error) {
	var strategyClients []*strategyClient
	config, err := b.getConfig()
	if err != nil {
		return nil, err
	}

	for _, v := range config.Clients {
		secure, err := strconv.ParseBool(v["secure"])
		if err != nil {
			return nil, err
		}

		endpoint := v["endpoint"]
		accessKeyID := v["access_key_id"]
		secretAccessKey := v["secret_access_key"]

		newMinioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: secure,
		})
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(v["weight"], 64)
		if err != nil {
			return nil, err
		}

		strategyClient := newStrategyClient(newMinioClient, weight)
		strategyClients = append(strategyClients, strategyClient)
	}
	return strategyClients, nil
}

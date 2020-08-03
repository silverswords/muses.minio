package storage

import (
	"strconv"

	"github.com/minio/minio-go/v6"
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

		newMinioClient, err := minio.New(v["endpoint"], v["accessKeyID"], v["secretAccessKey"], secure)
		if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(v["weight"], 64)
		if err != nil {
			return nil, err
		}

		strategyClient := newStrategyClient(newMinioClient, weight)
		if err != nil {
			return nil, err
		}

		strategyClients = append(strategyClients, strategyClient)
	}
	return strategyClients, nil
}

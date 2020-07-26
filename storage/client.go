package storage

import (
	"log"
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

func getStrategyClients() []*strategyClient {
	var strategyClients []*strategyClient
	for _, v := range getConfig().Clients {
		newMinioClient, err := minio.New(v["endpoint"], v["accessKeyID"], v["secretAccessKey"], true)
		if err != nil {
			log.Fatalln(err)
		}

		weight, err := strconv.ParseFloat(v["weight"], 64)
		if err != nil {
			log.Fatalln(err)
		}

		strategyClient := newStrategyClient(newMinioClient, weight)
		if err != nil {
			log.Fatalln(err)
		}

		strategyClients = append(strategyClients, strategyClient)
	}
	return strategyClients
}

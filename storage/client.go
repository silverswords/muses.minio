package storage

import (
	"log"
	"strconv"

	"github.com/minio/minio-go/v6"
)

// type Client interface {
// 	getMinioClient() *minio.Client
// }

type client *minio.Client

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

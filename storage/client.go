package storage

import (
	"log"

	"github.com/minio/minio-go/v6"
)

type Client interface {
	getMinioClient() *minio.Client
}

type client *minio.Client

func getMinioClients() []*minio.Client {
	// var client *minio.Client
	var minioClients []*minio.Client
	for _, v := range getConfig().Clients {
		// client := getConfig().Clients
		newMinioClient, err := minio.New(v["endpoint"], v["accessKeyID"], v["secretAccessKey"], true)
		if err != nil {
			log.Fatalln(err)
		}

		minioClients = append(minioClients, newMinioClient)
	}
	return minioClients
}

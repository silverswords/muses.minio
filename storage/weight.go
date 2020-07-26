package storage

import (
	"math/rand"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) saveByWeight() *minio.Client {
	var minioClient *minio.Client
	var weightflag float64
	random := rand.Float64()
	length := len(b.strategyClients)
	for i := 0; i < length; i++ {
		for j := 0; j < length-1-i; j++ {
			if b.strategyClients[j].weight > b.strategyClients[j+1].weight {
				b.strategyClients[j], b.strategyClients[j+1] = b.strategyClients[j+1], b.strategyClients[j]
			}
		}
	}

	for _, v := range b.strategyClients {
		weightflag += v.weight
		v.weight = weightflag
	}

	for _, v := range b.strategyClients {
		if random < v.weight {
			return v.client
		}
	}

	return minioClient
}

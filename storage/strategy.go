package storage

import (
	"math/rand"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) saveByWeight() *minio.Client {
	var minioClient *minio.Client
	// var weights []float64
	random := rand.Float64()
	for _, v := range b.strategyClients {
		if random < v.weight {
			return v.client
		}
	}
	// sort.Float64s(weights)
	// sort.Sort(sort.Reverse(sort.Float64Slice(weights)))

	return minioClient
}

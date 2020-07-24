package storage

import (
	"math/rand"
	"os"

	"github.com/minio/minio-go/v6"
)

type Weight struct {
	minioWeight map[*minio.Client]float64
}

func (b *Bucket) weightSave(weight float64) {
	b.minioWeight[b.minioClient] = weight
}

func (b *Bucket) weightGet() float64 {
	return b.minioWeight[b.minioClient]
}

func (b *Bucket) SaveBaseOnWeight(objectName string, object *os.File) {
	random := rand.Float64()
	var scale float64
	for _, v := range b.clients {
		scale += b.minioWeight[v]
		if random <= scale {
			b.PutObject(objectName, object)
		}
	}
}

package storage

import (
	"math/rand"
	"time"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) saveByWeight() (*minio.Client, error) {
	var weightflag float64
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	random := r.Float64()

	clients, err := b.getStrategyClients()
	if err != nil {
		return nil, err
	}

	length := len(clients)
	strategyClient := clients
	for i := 0; i < length; i++ {
		for j := 0; j < length-1-i; j++ {
			if strategyClient[j].weight > strategyClient[j+1].weight {
				strategyClient[j], strategyClient[j+1] = strategyClient[j+1], strategyClient[j]
			}
		}
	}

	for _, v := range strategyClient {
		weightflag += v.weight
		v.weight = weightflag
	}

	for _, v := range strategyClient {
		if random < v.weight {
			return v.client, nil
		}
	}

	return nil, nil
}

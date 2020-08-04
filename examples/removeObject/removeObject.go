package main

import (
	"log"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	b := storage.NewBucket("test", "config.yaml", "../", storage.WithStrategy("weightStrategy"))
	exists, err := b.CheckBucket("test")
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		err = b.MakeBucket()
		if err != nil {
			log.Println("errors in MakeBucket", err)
		}
	}

	err = b.RemoveObject("cat")
	if err != nil {
		log.Println("errors in RemoveObject", err)
	}
}

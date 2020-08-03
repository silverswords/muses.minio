package main

import (
	"log"
	"os"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	object, err := os.Open("../moon.jpg")
	if err != nil {
		log.Println("errors in Openfile", err)
	}
	defer object.Close()

	b := storage.NewBucket("test", "cn-north-1", "weightStrategy", "config.yaml", "../")
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

	err = b.PutObject("moon", object)
	if err != nil {
		log.Println("errors in PutObject", err)
	}
}

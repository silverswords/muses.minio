package main

import (
	"github.com/silverswords/muses.minio/storage"
	"log"
	"os"
)

func main() {
	object, err := os.Open("../moon.jpg")
	if err != nil {
		log.Println("errors in Openfile", err)
	}
	defer object.Close()

	b := storage.NewBucketConfig("test", "config.yaml", "../", storage.OtherOptions{})
	exists, err := b.CheckBucket()
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		err = b.MakeBucket()
		if err != nil {
			log.Println("errors in MakeBucket", err)
		}
	}

	err = b.PutObject("cat", object)
	if err != nil {
		log.Println("errors in PutObject", err)
	}
}

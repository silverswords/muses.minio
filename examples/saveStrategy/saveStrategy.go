package main

import (
	"github.com/silverswords/muses.minio/storage"
	"log"
	"os"
)

func main() {
	object, err := os.Open("../cat.jpg")
	if err != nil {
		log.Println("errors in Openfile", err)
	}
	defer object.Close()

	// savebyweight
	byWeight := storage.NewBucketConfig("apple","config.yaml", "../", storage.WithStrategy("weightStrategy"))

	exists, err := byWeight.CheckBucket()
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		err = byWeight.MakeBucket()
		if err != nil {
			log.Println("errors in MakeBucket", err)
		}
	}

	err = byWeight.PutObject("cat", object)
	if err != nil {
		log.Println("errors in PutObject", err)
	}

	// multiwrite
	multiwrite := storage.NewBucketConfig("banana", "config.yaml", "../")
	exists, err = multiwrite.CheckBucket()
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		err = multiwrite.MakeBucket()
		if err != nil {
			log.Println("errors in MakeBucket", err)
		}
	}

	err = multiwrite.PutObject("cat", object)
	if err != nil {
		log.Println("errors in PutObject", err)
	}
}

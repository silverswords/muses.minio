package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	object, err := os.Open("./cat.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()

	// savebyweight
	byWeight := storage.NewBucket("apple", "cn-north-1", "weightStrategy", "config.yaml", ".")

	exists, err := byWeight.CheckBucket("apple")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		err = byWeight.MakeBucket()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = byWeight.PutObject("cat", object)
	if err != nil {
		fmt.Println(err)
	}

	// multiwrite
	multiwrite := storage.NewBucket("banana", "cn-north-1", "multiWriteStrategy", "config.yaml", ".")
	exists, err = multiwrite.CheckBucket("banana")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		err = multiwrite.MakeBucket()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = multiwrite.PutObject("cat", object)
	if err != nil {
		fmt.Println(err)
	}
}

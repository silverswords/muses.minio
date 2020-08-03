package main

import (
	"fmt"
	"log"
	"os"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	object, err := os.Open("./moon.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	b := storage.NewBucket("test", "cn-north-1", "weightStrategy", "config.yaml", ".")
	exists, err := b.CheckBucket("test")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		err = b.MakeBucket()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = b.PutObject("moon", object)
	if err != nil {
		fmt.Println(err)
	}
}

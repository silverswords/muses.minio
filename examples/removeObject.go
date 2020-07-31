package main

import (
	"fmt"
	"log"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	b := storage.NewBucket("test", "cn-north-1", "weightStrategy")

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

	err = b.RemoveObject("cat")
	if err != nil {
		fmt.Println(err)
	}
}

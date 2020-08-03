package main

import (
	"fmt"
	"log"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	b := storage.NewBucket("test", "cn-north-1", "weightStrategy", "config.yaml", ".")
	err := b.MakeBucket()
	if err != nil {
		log.Fatalln(err)
	}
	exists, err := b.CheckBucket("test")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		fmt.Println("bucket does not exists.")
		// err = b.MakeBucket()
		// if err != nil {
		// 	log.Fatalln(err)
		// }
	}
}

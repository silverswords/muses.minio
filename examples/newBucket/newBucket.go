package main

import (
	"log"

	"github.com/silverswords/muses.minio/storage"
)

func main() {
	b := storage.NewBucket("test", "cn-north-1", "weightStrategy", "config.yaml", "../")
	err := b.MakeBucket()
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	exists, err := b.CheckBucket("test")
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		log.Println("bucket does not exists.")
		// err = b.MakeBucket()
		// if err != nil {
		// 	log.Println(err)
		// }
	}
}

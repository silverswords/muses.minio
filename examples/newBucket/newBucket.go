package main

import (
	"github.com/silverswords/muses.minio/storage"
	"log"
)

func main() {
	b := storage.NewBucketConfig("test", "config.yaml", "../", storage.OtherBucketConfigOptions{})
	err := b.MakeBucket()
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	exists, err := b.CheckBucket()
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

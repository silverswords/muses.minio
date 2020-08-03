package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/storage"
	"io"
	"log"
	"os"
)

func main() {
	b := storage.NewBucket("test", "cn-north-1", "weightStrategy", "config.yaml", "../")

	exists, err := b.CheckBucket("test")
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		//err = b.MakeBucket()
		//if err != nil {
		//	log.Println(err)
		//}
		log.Println("Bucket does not exist.")
	}

	minioObject, err := b.GetObject("moon")
	if err != nil {
		log.Println("errors in GetObject", err)
	}

	file, err := os.Create("file")
	var buffer = bytes.NewBuffer(minioObject)
	_, err = io.Copy(file, buffer)
	if err != nil {
		log.Println("errors in Copy", err)
	}
}

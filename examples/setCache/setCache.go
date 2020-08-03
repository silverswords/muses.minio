package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/storage"
	"io"
	"log"
	"os"
)

func main() {
	b := storage.NewBucket("banana", "cn-north-1", "weightStrategy", "config.yaml", "../")

	exists, err := b.CheckBucket("banana")
	if exists && err != nil {
		log.Println("errors in CheckBucket", err)
	}
	if !exists {
		log.Println("bucket does not exist.")
	}

	minioObject, err := b.GetObject("cat")
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

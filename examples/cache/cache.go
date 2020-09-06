package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"log"
	"os"
)

func main() {
	cb, err := bucketStorage.NewBucketWithCacheConfig("moon", "config.yaml", "../")
	if err != nil {
		log.Println("errors in NewBucketConfig", err)
	}

	err = cb.MakeBucket()
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	ok, err := cb.CheckBucket()
	if err != nil {
		log.Println("errors in CheckBucket", err)
	}

	file, err := os.Open("../bluemoon.jpg")
	if err != nil {
		log.Println("errors in Open", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Println("errors in Stat", err)
	}

	if ok {
		err = cb.PutObject("bluemoon", file, stat.Size())
		if err != nil {
			log.Println("errors in PutObject", err)
		}

		minioObject, err := cb.GetObject("bluemoon")
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
}

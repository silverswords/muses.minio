package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"log"
	"os"
)

func main() {
	cb, err := bucketStorage.NewBucket("moon", "config.yaml", "../")
	if err != nil {
		log.Println("errors in NewBucket", err)
	}

	err = bucketStorage.MakeBucket(cb)
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	ok, err := bucketStorage.CheckBucket(cb)
	if err != nil {
		log.Println("errors in CheckBucket", err)
	}

	file, err := os.Open("../bluemoon.jpg")
	if err != nil {
		log.Println("errors in Open", err)
	}
	defer file.Close()

	if ok {
		err = bucketStorage.PutObject(cb, "bluemoon", file)
		if err != nil {
			log.Println("errors in PutObject", err)
		}

		minioObject, err := bucketStorage.GetObject(cb,"bluemoon")
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

package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"log"
	"os"
)

func main() {
	bucket, err := bucketStorage.NewBucketConfig("test", "config.yaml", "./examples")
	if err != nil {
		log.Println("errors in NewBucketConfig", err)
	}

	err = bucket.MakeBucket()
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	ok, err := bucket.CheckBucket()
	if err != nil {
		log.Println("errors in CheckBucket", err)
	}

	object, err := os.Open("../cat.jpg")
	if err != nil {
		log.Println("errors in Open", err)
	}
	defer object.Close()

	if ok {
		err = bucket.PutObject("cat", object)
		if err != nil {
			log.Println("errors in PutObject", err)
		}

		minioObject, err := bucket.GetObject("moon")
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

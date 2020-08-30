package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"log"
	"os"
)

func main() {
	bucket, err := bucketStorage.NewBucketConfig("test", "config.yaml", "../")
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

	file, err := os.Open("../cat.jpg")
	if err != nil {
		log.Println("errors in Open", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Println(err)
	}

	if ok {
		err = bucket.PutObject("cat", file, stat.Size())
		if err != nil {
			log.Println("errors in PutObject", err)
		}

		minioObject, err := bucket.GetObject("cat")
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

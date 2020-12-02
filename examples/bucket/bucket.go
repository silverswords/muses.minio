package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/bucketStorage"
	"io"
	"log"
	"os"
)

func main() {
	bucket, err := bucketStorage.NewBucket("test", "config.yaml", "../")
	if err != nil {
		log.Println("errors in NewBucketConfig", err)
	}

	err = bucketStorage.MakeBucket(bucket)
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	ok, err := bucketStorage.CheckBucket(bucket)
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
		log.Println("errors in Stat", err)
	}

	if ok {
		err = bucketStorage.PutObject(bucket,"cat", file, bucketStorage.WithObjectSize(stat.Size()))
		if err != nil {
			log.Println("errors in PutObject", err)
		}

		minioObject, err := bucketStorage.GetObject(bucket,"cat")
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

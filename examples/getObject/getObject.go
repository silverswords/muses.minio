package main

import (
	"bytes"
	"github.com/silverswords/muses.minio/storage"
	"io"
	"log"
	"os"
)

func main() {
	b := storage.NewBucketConfig("test", "config.yaml", "../", storage.OtherBucketConfigOptions{})

	exists, err := b.CheckBucket()
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

	minioObject, err := b.GetObject("cat", storage.ObjectServerSideEncryptions{})
	if err != nil {
		log.Println("errors in GetObject", err)
	}

	file, err := os.Create("testfile")
	var buffer = bytes.NewBuffer(minioObject)
	_, err = io.Copy(file, buffer)
	if err != nil {
		log.Println("errors in Copy", err)
	}
}

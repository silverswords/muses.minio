package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func checkBucket(bucketName string) error {
	location := "cn-north-1"
	new := NewMinioClient()
	minioClient := new.GetMinioClient()
	exists, err := minioClient.BucketExists(bucketName)
	if exists == false && err == nil {
		err = minioClient.MakeBucket(bucketName, location)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return err
}

func Save(bucketName string, filePath string) (int64, error) {
	new := NewMinioClient()
	minioClient := new.GetMinioClient()

	object, err := os.Open(bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	err = checkBucket(bucketName)
	if err != nil {
		log.Fatalln(err)
		return 0, err
	} else {
		n, err := minioClient.PutObject(bucketName, filePath, object, objectStat.Size(), minio.PutObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return n, nil
	}
}

package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func SaveMinioObject(bucketName string, objectName string) error {
	new, err := NewMinioClient()
	if err != nil {
		log.Fatalln(err)
	}
	minioClient := new.GetMinioClient()

	reader, err := os.Open(bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()
	objectStat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	err = checkBucket(bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = minioClient.PutObject(bucketName, objectName, reader, objectStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	var b *bucketObjectCache
	err = b.Set(bucketName, objectName)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func Save(bucketName string, objectName string) (int64, error) {
	new, err := NewMinioClient()
	if err != nil {
		return 0, err
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
		return 0, err
	}

	n, err := minioClient.PutObject(bucketName, objectName, reader, objectStat.Size(), minio.PutObjectOptions{})
	if err != nil {
		return 0, err
	}

	var b *bucketObjectCache
	err = b.Set(bucketName, objectName)
	if err != nil {
		return 0, err
	}

	return n, nil
}

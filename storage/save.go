package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func (m *minioClient) Save(bucketName string, objectName string) error {
	reader, err := os.Open(bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()
	objectStat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	mc := Open("minio://minio:minio123@127.0.0.1:9001")
	exists, err := mc.CheckBucket(bucketName)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	minioClient := m.newMinioClient
	if exists {
		_, err = minioClient.PutObject(bucketName, objectName, reader, objectStat.Size(), minio.PutObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		err = objectCache.Set(bucketName, objectName)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("Bucket does not exist.")
	}

	return err
}

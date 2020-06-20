package storage

import (
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
			return err
		}
	}

	return err
}

func Save(bucketName string, objectName string) (int64, error) {
	new := NewMinioClient()
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

func Get(bucketName string, objectName string) (*minio.Object, error) {
	new := NewMinioClient()
	minioClient := new.GetMinioClient()

	var b *bucketObjectCache
	minioObject := b.Get(bucketName, objectName)
	if minioObject == nil {
		minioObject, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		return minioObject, nil
	}

	return minioObject, nil
}

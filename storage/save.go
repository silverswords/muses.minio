package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) Save(objectName string) error {
	reader, err := os.Open(b.bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()
	objectStat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	exists, err := b.checkBucket()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	if exists {
		_, err = b.client.getMinioClient().PutObject(b.bucketName, objectName, reader, objectStat.Size(), minio.PutObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		err = b.cacheSave(objectName)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("Bucket does not exist.")
	}

	return err
}

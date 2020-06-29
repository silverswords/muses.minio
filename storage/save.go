package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) Save(objectName string, object *os.File) error {
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	exists, err := b.checkBucket()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	if exists {
		_, err = b.client.getMinioClient().PutObject(b.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
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

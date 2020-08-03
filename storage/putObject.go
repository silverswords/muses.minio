package storage

import (
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

func (b *Bucket) PutObject(objectName string, object *os.File) error {
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	exists, err := b.CheckBucket(b.bucketName)
	if err != nil {
		return err
	}

	var buf = make([]byte, objectStat.Size())
	if exists {
		err = b.setCacheObject(buf, objectName)
		if err != nil {
			log.Println("errors in set cache: ",err)
		}
	}

	if b.strategy == "weightStrategy" {
		c, err := b.saveByWeight()
		if err != nil {
			return err
		}

		if exists {
			_, err = c.PutObject(b.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
			if err != nil {
				return err
			}
		} else {
			log.Println("Bucket does not exist.")
		}

	}
	if b.strategy == "multiWriteStrategy" {
		clients, err := b.getStrategyClients()
		if err != nil {
			return err
		}

		for _, v := range clients {
			if exists {
				_, err = v.client.PutObject(b.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
				if err != nil {
					return err
				}
			} else {
				log.Println("Bucket does not exist.")
			}
		}
	}

	//if exists {
	//	buf, err := b.GetObject(objectName)
	//	if err != nil {
	//		return err
	//	}
	//	err = b.setCacheObject(buf, objectName)
	//	if err != nil {
	//		log.Println("errors in set cache: ",err)
	//	}
	//}

	return nil
}

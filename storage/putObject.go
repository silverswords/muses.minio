package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

type ObjectEncryptions struct {
	encryption bool
	password string
}

func (b *Bucket) PutObject(objectName string, object *os.File, opts ObjectEncryptions) error {
	var e encrypt.ServerSide
	if opts.encryption && opts.password != "" {
		e = encrypt.DefaultPBKDF([]byte(opts.password), []byte(b.bucketName + objectName))
	}

	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	exists, err := b.CheckBucket()
	if err != nil {
		return err
	}

	var buf = make([]byte, objectStat.Size())
	if b.cache && exists {
		err = b.setCacheObject(buf, objectName)
		if err != nil {
			log.Println("errors in set cache: ",err)
		}
	}

	if b.Strategy == "weightStrategy" {
		c, err := b.saveByWeight()
		if err != nil {
			return err
		}

		if exists {
			_, err = c.PutObject(context.Background(), b.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: e})
			if err != nil {
				return err
			}
		} else {
			log.Println("Bucket does not exist.")
		}
	}

	if b.Strategy == "" || b.Strategy == "multiWriteStrategy" {
		clients, err := b.getStrategyClients()
		if err != nil {
			return err
		}

		for _, v := range clients {
			if exists {
				_, err = v.client.PutObject(context.Background(), b.bucketName, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: e})
				if err != nil {
					return err
				}
			} else {
				log.Println("Bucket does not exist.")
			}
		}
	}

	return nil
}

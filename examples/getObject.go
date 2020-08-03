package main

import (
	"fmt"
	"github.com/silverswords/muses.minio/storage"
	"io"
	"log"
	"os"
)

func main() {
	b := storage.NewBucket("banana", "cn-north-1", "weightStrategy", "config.yaml", ".")

	exists, err := b.CheckBucket("banana")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		err = b.MakeBucket()
		if err != nil {
			log.Fatalln(err)
		}
	}

	minioObject, err := b.GetObject("cat")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(minioObject)

	//localFile, err := os.Create("catfile")
	//_, err = io.Copy(localFile, minioObject)
	//if err != nil {
	//	log.Println(err)
	//}
}

package main

import (
	"fmt"
	"github.com/silverswords/muses.minio/storage"
	"io"
	"log"
)

func main() {
	b := storage.NewBucket("banana", "cn-north-1", "weightStrategy")

	exists, err := b.CheckBucket("banana")
	if exists && err != nil {
		log.Fatalln(err)
	}
	if !exists {
		fmt.Println("bucket does not exist.")
	}

	minioObject, err := b.GetObject("cat")
	if err != nil {
		fmt.Println(err)
	}

	stat, err := minioObject.Stat()
	var buf = make([]byte, stat.Size)
	_, err = io.ReadFull(minioObject, buf)
	if err != nil {
		log.Println(err)
	}

	//ctx := context.TODO()
	//storage.SetCacheObject("cat", buf, ctx)
	//objectByte := storage.GetCache("cat", ctx)
	//fmt.Println(objectByte)
	//
	//var writer bytes.Buffer
	//
	//_, err = writer.Write(objectByte)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//fmt.Println(writer.String())
}

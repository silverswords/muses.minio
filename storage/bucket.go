package storage

import "log"

type Bucket struct {
	bucketName string
	location   string
}

func checkBucket(bucketName string) (bool, error) {
	exists, err := m.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

func NewBucket(bucketName string) *Bucket {
	location := "cn-north-1"
	err := m.MakeBucket(bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	return &Bucket{
		bucketName: bucketName,
		location:   location,
	}
}

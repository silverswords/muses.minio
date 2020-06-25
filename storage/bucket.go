package storage

import "log"

type Bucket struct {
	BucketName string
	Location   string
}

func checkBucket(bucketName string) (bool, error) {
	exists, err := m.BucketExists(bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

func newBucket(bucketName string) *Bucket {
	location := "cn-north-1"
	err := m.MakeBucket(bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	return &Bucket{
		BucketName: bucketName,
		Location:   location,
	}
}

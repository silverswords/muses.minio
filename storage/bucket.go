package storage

import (
	"io/ioutil"
	"log"

	"github.com/minio/minio-go/v6"
	"gopkg.in/yaml.v2"
)

type Bucket struct {
	bucketName        string `yaml:"bucketName"`
	client            `yaml:"client"`
	bucketObjectCache `yaml:"bucketObjectCache"`
}

func newBucket(s string, bucketName string) *Bucket {
	return &Bucket{
		bucketName: bucketName,
		client: client{
			minioClient: minioClient{
				url: s,
			},
		},
		bucketObjectCache: bucketObjectCache{
			items: make(map[string]*minio.Object),
		},
	}
}

func (b *Bucket) makeBucket() error {
	location := "cn-north-1"
	err := b.client.getMinioClient().MakeBucket(b.bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	return err
}

func (b *Bucket) checkBucket() (bool, error) {
	exists, err := b.client.getMinioClient().BucketExists(b.bucketName)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	return exists, err
}

func (b *Bucket) GetConf() *Bucket {
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(yamlFile, b)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

package main

import (
	"encoding/xml"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/silverswords/muses.minio/bucketStorage"
	"log"
)

func main() {
	bucket, err := bucketStorage.NewBucket("test", "config.yaml", "../")
	if err != nil {
		log.Println("errors in NewBucketConfig", err)
	}

	err = bucket.MakeBucket()
	if err != nil {
		log.Println("errors in MakeBucket", err)
	}

	ok, err := bucket.CheckBucket()
	if err != nil {
		log.Println("errors in CheckBucket", err)
	}

	if ok {
		err := bucket.SetBucketVersioning()
		if err != nil {
			log.Println("errors in SetBucketVersioning", err)
		}

		replicationStr := `<ReplicationConfiguration><Rule><ID>rule1</ID><Status>Enabled</Status><Priority>1</Priority><DeleteMarkerReplication><Status>Disabled</Status></DeleteMarkerReplication><Destination><Bucket>arn:aws:s3:::test</Bucket></Destination><Filter><And><Prefix></Prefix></And></Filter></Rule></ReplicationConfiguration>`
		var replCfg replication.Config
		err = xml.Unmarshal([]byte(replicationStr), &replCfg)
		if err != nil {
			log.Println("errors in Unmarshal", err)
		}

		// This replication ARN should have been generated for replication endpoint using `mc admin bucket remote` command
		replCfg.Role = "arn:minio:replication:us-east-1:9f6c650d-6b5a-4e12-8e47-33f135dcb90f:test"

		err = bucket.SetBucketReplication(replCfg)
		if err != nil {
			log.Println("errors in SetBucketReplication", err)
		}
	}
}

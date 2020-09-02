package main

import (
	"encoding/xml"
	"github.com/minio/minio-go/v7/pkg/replication"
	"github.com/silverswords/muses.minio/bucketStorage"
	"log"
)

func main() {
	bucket, err := bucketStorage.NewBucketConfig("test", "config.yaml", "../")
	if err != nil {
		log.Println("errors in NewBucketConfig", err)
	}

	replicationStr := `<ReplicationConfiguration><Rule><ID>rule1</ID><Status>Enabled</Status><Priority>1</Priority><DeleteMarkerReplication><Status>Disabled</Status></DeleteMarkerReplication><Destination><Bucket>arn:aws:s3:::moon</Bucket></Destination><Filter><And><Prefix></Prefix><Tag><Key></Key><Value></Value></Tag><Tag><Key></Key><Value></Value></Tag></And></Filter></Rule></ReplicationConfiguration>`
	var replCfg replication.Config
	err = xml.Unmarshal([]byte(replicationStr), &replCfg)
	if err != nil {
		log.Fatalln(err)
	}

	replCfg.Role = ""
	err = bucket.SetBucketReplication(replCfg)
	if err != nil {
		log.Println("errors in SetBucketReplication", err)
	}
}

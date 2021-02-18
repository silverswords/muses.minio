package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silverswords/muses.minio/bucketStorage"
	object "github.com/silverswords/muses.minio/controller/gin"
	"log"
)

func main() {
	router := gin.Default()
	bucket, err := bucketStorage.NewBucket("ashe", "config.yaml", "../examples")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("hhh", bucket)
	con := object.New(bucket)
	con.RegisterRouter(router.Group("/api/v1/object"))

	_ = router.Run(":8001")
}

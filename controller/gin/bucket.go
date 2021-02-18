/*
 * Revision History:
 *     Initial: 2020/12/2      oiar
 */

package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silverswords/muses.minio/bucketStorage"
	"github.com/silverswords/muses.minio/bucketStorage/driver"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// BucketController -
type BucketController struct {
	bucket *bucketStorage.Bucket
}

// New -
func New(bucket *bucketStorage.Bucket) *BucketController {
	return &BucketController{
		bucket,
	}
}

// RegisterRouter -
func (b *BucketController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := b.bucket.MakeBucket()
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/upload", b.upload)
	r.POST("/delete", b.delete)
	r.POST("/download", b.download)
	r.POST("/list", b.listObjects)
}

func (b *BucketController) upload(c *gin.Context) {
	var (
		req struct {
			File multipart.FileHeader `json:"file"       binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	file, err := req.File.Open()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway,"error": err})
		return
	}

	fileSize := req.File.Size
	_, err = b.bucket.NewTypedWriter(context.Background(), req.File.Filename, file, driver.WithObjectSize(fileSize))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway,"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (b *BucketController) delete(c *gin.Context) {
	var (
		req struct {
			ObjectName  string `json:"objectName"      binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	fmt.Println(req.ObjectName, "------objectName-----")
	err = b.bucket.Delete(context.Background(), req.ObjectName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (b *BucketController) download(c *gin.Context) {
	var (
		req struct {
			ObjectName string `json:"objectName"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	u, err := b.bucket.SignedURL(context.Background(), req.ObjectName, time.Second*24*60*60, "GET")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	fmt.Println("url:", u)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "url": u})
}

func (b *BucketController) listObjects(c *gin.Context) {
	var (
		req struct {
			ProjectID int `json:"ProjectId"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ch := b.bucket.ListObjects()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ch": ch})
}
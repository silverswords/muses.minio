/*
 * Revision History:
 *     Initial: 2020/12/2      oiar
 */

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silverswords/muses.minio/bucketStorage"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
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

	err := bucketStorage.MakeBucket(b.bucket)
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
			ObjectName  string `json:"objectName"      binding:"required"`
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
	err = bucketStorage.PutObject(b.bucket, req.ObjectName, file, bucketStorage.WithObjectSize(fileSize))
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
	err = bucketStorage.RemoveObject(b.bucket, req.ObjectName)
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

	//object, err := bucketStorage.GetObject(b.bucket, req.ObjectName)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"file\"")
	u, err := bucketStorage.PresignedGetObject(b.bucket, req.ObjectName, time.Second*24*60*60, reqParams)
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

	ch := bucketStorage.ListObjects(b.bucket)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ch": ch})
}
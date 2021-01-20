package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CreateBucketAPI interface {
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, opts ...func(options *s3.Options)) (*s3.CreateBucketOutput, error)
}

func MakeBucket(ctx context.Context, api CreateBucketAPI, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return api.CreateBucket(ctx, input)
}

func main() {
	bucket := flag.String("oiar", "oiar", "The name of the bucket")
	flag.Parse()

	if *bucket == "" {
		fmt.Println("You must supply a bucket name (-b BUCKET)")
		return
	}

	accessKey := "snowoiar@gmail.com"
	secretKey := "Lipsgo"

	config, err := config.LoadDefaultConfig(context.TODO())
	c := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	if err != nil {
		fmt.Println(err)
	}

	client := s3.NewFromConfig(config, func(o *s3.Options) {
		o.Region = "us-west-2"
		o.Credentials = c
	})

	input := &s3.CreateBucketInput{
		Bucket: bucket,
	}

	_, err = MakeBucket(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err, "Cannot create bucket" + *bucket)
	}
}
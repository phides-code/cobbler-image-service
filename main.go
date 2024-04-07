package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

var myS3 BucketBasics

func init() {
	s3Client, err := getClient()

	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of BucketBasics with the S3 client
	myS3 = BucketBasics{
		S3Client: &s3Client,
	}
}

func main() {
	lambda.Start(router)
}

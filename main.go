package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var myS3 s3.Client

func init() {
	s3Client, err := getClient()

	if err != nil {
		log.Fatal(err)
	}

	myS3 = s3Client
}

func main() {
	lambda.Start(router)
}

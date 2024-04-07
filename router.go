package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
)

type ResponseStructure struct {
	Data         interface{} `json:"data"`
	ErrorMessage *string     `json:"errorMessage"` // can be string or nil
}

var validate *validator.Validate = validator.New()

var headers = map[string]string{
	"Access-Control-Allow-Origin":  OriginURL,
	"Access-Control-Allow-Headers": "Content-Type",
}

type ImageData struct {
	Image   string `json:"image"`
	FileExt string `json:"fileExt"`
}

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return processPost(req)
	case "OPTIONS":
		return processOptions()
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func processOptions() (events.APIGatewayProxyResponse, error) {
	additionalHeaders := map[string]string{
		"Access-Control-Allow-Methods": "OPTIONS, POST",
		"Access-Control-Max-Age":       "3600",
	}
	mergedHeaders := mergeHeaders(headers, additionalHeaders)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    mergedHeaders,
	}, nil
}

func processPost(
	req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var imageData ImageData

	err := json.Unmarshal([]byte(req.Body), &imageData)

	if err != nil {
		log.Println("Error decoding request body:", err)
		return clientError(http.StatusBadRequest)
	}

	err = validate.Struct(&imageData)

	if err != nil {
		log.Println("Error decoding request body:", err)
		return clientError(http.StatusBadRequest)
	}

	// Decode the base64-encoded image data
	imageBytes, err := base64.StdEncoding.DecodeString(imageData.Image)
	if err != nil {
		log.Println("Error decoding base64 image:", err)
		return clientError(http.StatusBadRequest)
	}

	contentType := http.DetectContentType(imageBytes)

	// Check if the uploaded file is an image
	if !strings.HasPrefix(contentType, "image/") {
		return clientError(http.StatusBadRequest)
	}

	fileExt := imageData.FileExt
	image := bytes.NewReader(imageBytes)

	// Create an instance of BucketBasics with the S3 client
	basics := BucketBasics{
		S3Client: &myS3,
	}

	fileName, err := basics.UploadFile(image, fileExt, contentType)

	if err != nil {
		return serverError(err)
	}

	response := ResponseStructure{
		Data:         fileName,
		ErrorMessage: nil,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(responseJson),
		Headers:    headers,
	}, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func clientError(status int) (events.APIGatewayProxyResponse, error) {

	errorString := http.StatusText(status)

	response := ResponseStructure{
		Data:         nil,
		ErrorMessage: &errorString,
	}

	responseJson, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseJson),
		StatusCode: status,
		Headers:    headers,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err.Error())

	errorString := http.StatusText(http.StatusInternalServerError)

	response := ResponseStructure{
		Data:         nil,
		ErrorMessage: &errorString,
	}

	responseJson, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseJson),
		StatusCode: http.StatusInternalServerError,
		Headers:    headers,
	}, nil
}

func mergeHeaders(baseHeaders, additionalHeaders map[string]string) map[string]string {
	mergedHeaders := make(map[string]string)
	for key, value := range baseHeaders {
		mergedHeaders[key] = value
	}
	for key, value := range additionalHeaders {
		mergedHeaders[key] = value
	}
	return mergedHeaders
}

// Function to check if the provided bytes represent a valid image
func isImage(data []byte) bool {
	// Open the image using image.DecodeConfig to check the format
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		log.Println("Error decoding image:", err)
		return false
	}

	// Check if the format is one of the supported image formats
	supportedFormats := map[string]bool{
		"jpeg": true,
		"png":  true,
		"gif":  true,
	}
	return supportedFormats[strings.ToLower(format)]
}

package main

import (
	"log"
	"io"
	"net/http"
	"golang.org/x/net/html/charset"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(HandleRequest)
}

type LambdaEvent struct {
	Name string `json:"name"`
}

type LambdaResponse struct {
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, name LambdaEvent) (LambdaResponse, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type") // Optional, better guessing
	log.Printf("ContentType is %s", contentType)
    utf8reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(utf8reader)
	if err != nil {
		// handle error
		log.Fatal(err)
	}

	log.Println(string(body))

	return LambdaResponse{
		Message: "OK " + name.Name,
	}, nil
}
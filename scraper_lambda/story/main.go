package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"context"
	"github.com/rivernews/GoTools"
	"fmt"
)


func main() {
	lambda.Start(HandleRequest)
}

type LambdaResponse struct {
	OK bool `json:"OK:"`
	Message string `json:"message:"`
}

// SQS event
// refer to https://github.com/aws/aws-lambda-go/blob/v1.26.0/events/README_SQS.md
func HandleRequest(ctx context.Context, event events.SQSEvent) (LambdaResponse, error) {
	for _, message := range event.Records {
		GoTools.Logger("INFO", fmt.Sprintf("Story consumer! story URL: %s", message.Body))

		// TODO: fetch and archive
	}

	return LambdaResponse{
		OK: true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

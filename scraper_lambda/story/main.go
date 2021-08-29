package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"context"
	"github.com/rivernews/GoTools"
	"fmt"
)


func main() {
	lambda.Start(HandleRequest)
}

type LambdaEvent struct {
	StoryURL string `json:"storyURL"`
}

type LambdaResponse struct {
	OK bool `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	// TODO

	GoTools.SendSlackMessage(fmt.Sprintf("Story consumer! story URL: %s", event.StoryURL))

	return LambdaResponse{
		OK: true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

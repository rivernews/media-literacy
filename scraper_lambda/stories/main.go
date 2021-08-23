// You must use `main` package for lambda
// https://stackoverflow.com/a/50701572/9814131
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
	LandingURL string `json:"landingURL"`
}

type LambdaResponse struct {
	OK bool `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	// TODO

	GoTools.SendSlackMessage(fmt.Sprintf("Batch stories lambda started! Landing page S3 path: %s", event.LandingURL))

	return LambdaResponse{
		OK: true,
		Message: "Batch story fetch parse ok",
	}, nil
}

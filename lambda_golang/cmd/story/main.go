package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"context"

	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/newssite"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	lambda.Start(HandleRequest)
}

type LambdaResponse struct {
	OK      bool   `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, stepFunctionMapIterationInput newssite.StepFunctionMapIterationInput) (LambdaResponse, error) {
	GoTools.Logger("INFO", "Fetch single story launched")

	GoTools.Logger("INFO", fmt.Sprintf("%s", stepFunctionMapIterationInput.Story.Name))

	return LambdaResponse{
		OK:      true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

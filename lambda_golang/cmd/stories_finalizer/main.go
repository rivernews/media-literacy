package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"context"

	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/newssite"
)

func main() {
	lambda.Start(HandleRequest)
}

type LambdaResponse struct {
	OK      bool   `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, stepFunctionInput newssite.StepFunctionInput) (LambdaResponse, error) {
	GoTools.Logger("INFO", "Stories finalizer launched")

	items := newssite.DynamoDBQueryByS3Key(ctx, stepFunctionInput.LandingPageS3Key)

	if len(*items) != 1 {
		GoTools.Logger("ERROR", fmt.Sprintf(
			"Stories finalizer expect exactly one landing page `%s`, but query resulted in `%d` items",
			stepFunctionInput.LandingPageS3Key,
			len(*items),
		))
	}

	newssite.DynamoDBUpdateItemAddEvent(ctx, stepFunctionInput.LandingPageUuid, newssite.GetEventLandingStoriesFetched(stepFunctionInput.LandingPageS3Key))

	return LambdaResponse{
		OK:      true,
		Message: "OK",
	}, nil
}

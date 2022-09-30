package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"context"
	"math/rand"
	"time"

	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/cloud"
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

var sfnClient *sfn.Client

// SQS event
// refer to https://github.com/aws/aws-lambda-go/blob/v1.26.0/events/README_SQS.md
func HandleRequest(ctx context.Context, S3Event events.S3Event) (LambdaResponse, error) {
	GoTools.Logger("INFO", "Landing page S3 trigger launched")

	for _, record := range S3Event.Records {
		GoTools.Logger("INFO", fmt.Sprintf("S3 event ``` %s ```\n ", GoTools.AsJson(record)))
		landingPageS3Key := record.S3.Object.URLDecodedKey
		landingPageS3KeyTokens := strings.Split(landingPageS3Key, "/")
		newsSiteAlias := landingPageS3KeyTokens[0]
		landingPageCreatedAt := landingPageS3KeyTokens[len(landingPageS3KeyTokens)-2]

		out := cloud.DynamoDBPutItem(ctx, newssite.MediaTableItem{
			CreatedAt: landingPageCreatedAt,
			S3Key:     landingPageS3Key,
			DocType:   newssite.DOCTYPE_LANDING,
			Events: []newssite.MediaTableItemEvent{
				newssite.GetEventLandingPageFetched(newsSiteAlias, landingPageS3Key),
				newssite.GetEventLandingMetadataRequested(landingPageS3Key),
			},
			IsDocTypeWaitingForMetadata: newssite.DOCTYPE_LANDING,
		})
		GoTools.Logger("DEBUG", fmt.Sprintf("```%s```\n", GoTools.AsJson(out)))
	}

	return LambdaResponse{
		OK:      true,
		Message: "OK",
	}, nil
}

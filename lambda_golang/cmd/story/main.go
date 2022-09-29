package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"context"

	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/cloud"
	"github.com/rivernews/media-literacy/pkg/common"
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

	baseWait := 4
	waitRange := 100
	totalWait := rand.Intn(waitRange) + baseWait
	time.Sleep(time.Duration(totalWait) * time.Second)

	responseBody, _, _ := GoTools.Fetch(GoTools.FetchOption{
		URL: "https://checkip.amazonaws.com",
		QueryParams: map[string]string{
			"format": "json",
		},
		Method: "GET",
	})

	GoTools.Logger("INFO", fmt.Sprintf("IP=`%s` waited %d - %s", bytes.TrimSpace(responseBody), totalWait, stepFunctionMapIterationInput.Story.Name))
	storyS3Key := fmt.Sprintf("%s/stories/%s-%s/story.html", stepFunctionMapIterationInput.NewsSiteAlias, stepFunctionMapIterationInput.LandingPageTimeStamp, stepFunctionMapIterationInput.Story.Name)
	storyHtmlBodyText := common.Fetch(stepFunctionMapIterationInput.Story.URL)
	cloud.Archive(cloud.ArchiveArgs{
		BodyText: storyHtmlBodyText,
		Key:      storyS3Key,
	})

	cloud.DynamoDBPutItem(
		ctx,
		newssite.MediaTableItem{
			S3Key:   storyS3Key,
			DocType: newssite.DOCTYPE_STORY,
			Events: []newssite.MediaTableItemEvent{
				newssite.GetEventStoryFetched(
					stepFunctionMapIterationInput.Story.Name,
					stepFunctionMapIterationInput.Story.URL,
				),
			},
			IsDocTypeWaitingForMetadata: newssite.DOCTYPE_STORY,
		},
	)

	return LambdaResponse{
		OK:      true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

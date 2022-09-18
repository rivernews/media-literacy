package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
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
	GoTools.Logger("INFO", "Fetch stories lambda launched ... triggered by metadata.json creation.")

	for _, record := range S3Event.Records {

		GoTools.Logger("INFO", fmt.Sprintf("S3 event ``` %s ```\n ", GoTools.AsJson(record)))

		metadataS3KeyTokens := strings.Split(record.S3.Object.URLDecodedKey, "/")
		newsSiteAlias := metadataS3KeyTokens[0]
		landingPageTimeStamp := metadataS3KeyTokens[len(metadataS3KeyTokens)-2]

		metadataJSONString := cloud.Pull(record.S3.Object.URLDecodedKey)
		var metadata newssite.LandingPageMetadata
		GoTools.FromJson([]byte(metadataJSONString), &metadata)

		GoTools.Logger("INFO", fmt.Sprintf("Test first story: %d:%d", len(metadata.Stories), len(metadata.UntitledStories)))

		// fire step function, input = {stories: [{}, {}, {}...], newsSiteAlias:"", landingPageTimeStamp:""}
		sfnInput := GoTools.AsJson(&newssite.StepFunctionInput{
			Stories:              metadata.Stories,
			NewsSiteAlias:        newsSiteAlias,
			LandingPageTimeStamp: landingPageTimeStamp,
		})
		executionName := strings.ReplaceAll(fmt.Sprintf("%s--%s", landingPageTimeStamp, time.Now().Format(time.RFC3339)), ":", "")
		sfnArn := GoTools.GetEnvVarHelper("SFN_ARN")
		GoTools.Logger("INFO", fmt.Sprintf("execName: `%s`\nSfnArn: `%s`\nInput:``` %s ```\n", executionName, sfnArn, sfnInput))

		sfnOutput := cloud.SfnStartExecution(
			ctx,
			&sfn.StartExecutionInput{
				Input:           aws.String(sfnInput),
				Name:            aws.String(executionName),
				StateMachineArn: aws.String(sfnArn),
			},
		)

		GoTools.Logger("INFO", fmt.Sprintf("Sfn output ``` %s ```\n", GoTools.AsJson(sfnOutput)))

		/*
			storyChunk := message.Body
			GoTools.Logger("INFO", fmt.Sprintf("Story consumer! story chunk: %s", storyChunk))

			// TODO: fetch and archive for the chunk of storyURLs
			storyURLs := strings.Split(storyChunk, " ")

			for _, storyURL := range storyURLs {
				// TODO: randomized interval
				time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)

				storyS3Path := fmt.Sprintf("story/%s", storyURL)

				if !cloud.IsDuplicated(storyS3Path) {
					GoTools.Logger("INFO", fmt.Sprintf("Archiving story %s", storyURL))
				} else {
					GoTools.Logger("DEBUG", fmt.Sprintf("Skip archiving story %s", storyURL))
				}

				// _, resMessage, err := GoTools.Fetch(GoTools.FetchOption{
				// 	Method: "GET",
				// 	URL: "https://ipv4bot.whatismyipaddress.com",
				// })
				// if err != nil {
				// 	GoTools.Logger("ERROR", err.Error())
				// }
				// GoTools.Logger("INFO", fmt.Sprintf("%s-%s ip: %s", storyChunkId, storyURL, resMessage))
			}

		*/
	}

	return LambdaResponse{
		OK:      true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

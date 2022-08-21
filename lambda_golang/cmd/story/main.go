package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"context"
	"fmt"
	"strings"
	"time"
	"math/rand"

	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/cloud"
)


func main() {
	rand.Seed(time.Now().UnixNano())
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
		storyChunk := message.Body
		GoTools.Logger("INFO", fmt.Sprintf("Story consumer! story chunk: %s", storyChunk))

		// TODO: fetch and archive for the chunk of storyURLs
		storyURLs := strings.Split(storyChunk, " ")

		for _, storyURL := range storyURLs {
			// TODO: randomized interval
			time.Sleep(time.Duration(rand.Intn(5) + 5) * time.Second)

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
	}

	return LambdaResponse{
		OK: true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

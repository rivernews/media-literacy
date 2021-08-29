package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"context"
	"github.com/rivernews/GoTools"
	"fmt"

	"time"
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
		storyChunkId := message.Body
		GoTools.Logger("INFO", fmt.Sprintf("Story consumer! story chunk: %s", storyChunkId))

		// TODO: fetch and archive for the chunk of storyURLs

		for _, storyURL := range []string{ "http://story-00.com","http://story-01.com", } {
			// TODO: randomized interval
			time.Sleep(5 * time.Second)

			_, resMessage, err := GoTools.Fetch(GoTools.FetchOption{
				Method: "GET",
				URL: "https://ipv4bot.whatismyipaddress.com",
			})
			if err != nil {
				GoTools.Logger("ERROR", err.Error())
			}
			GoTools.Logger("INFO", fmt.Sprintf("%s-%s ip: %s", storyChunkId, storyURL, resMessage))
		}
	}

	return LambdaResponse{
		OK: true,
		Message: "Story consumer fetch parse ok",
	}, nil
}

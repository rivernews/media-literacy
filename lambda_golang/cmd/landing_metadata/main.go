// You must use `main` package for lambda
// https://stackoverflow.com/a/50701572/9814131
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	//"math"

	"github.com/aws/aws-lambda-go/events" // https://github.com/aws/aws-lambda-go/blob/main/events/README_S3.md
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rivernews/GoTools"

	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go-v2/config"
	//"github.com/aws/aws-sdk-go-v2/service/sqs"

	// local packages
	"github.com/rivernews/media-literacy/pkg/cloud"
	"github.com/rivernews/media-literacy/pkg/newssite"
)

func main() {
	lambda.Start(HandleRequest)
}

type StepFunctionInput struct {
	LandingS3Key string `json:"landingS3Key"`
}

type LambdaResponse struct {
	OK      bool   `json:"OK:"`
	Message string `json:"message:"`
}

type LandingPageMetadata struct {
	Stories         []newssite.Topic `json:"stories"`
	UntitledStories []newssite.Topic `json:"untitledstories"`
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) (LambdaResponse, error) {
	GoTools.Logger("INFO", "Landing page metadata.json generator launched")

	landingPageS3Key := s3Event.Records[0].S3.Object.URLDecodedKey
	GoTools.Logger("INFO", fmt.Sprintf("Captured landing page at %s", landingPageS3Key))

	landingPageHtmlText := cloud.Pull(landingPageS3Key)
	landingPageS3KeyTokens := strings.Split(landingPageS3Key, "/")
	metadataS3DirKeyTokens := landingPageS3KeyTokens[:len(landingPageS3KeyTokens)-1]
	metadataS3Key := fmt.Sprintf("%s/metadata.json", strings.Join(metadataS3DirKeyTokens, "/"))

	result := newssite.GetStoriesFromEconomy(landingPageHtmlText)

	metadataJSONBytes, metadataJSONStringError := json.Marshal(LandingPageMetadata{Stories: result.Topics, UntitledStories: result.UntitledTopics})
	if metadataJSONStringError != nil {
		GoTools.Logger("ERROR", metadataJSONStringError.Error())
	}
	metadataJSONString := string(metadataJSONBytes)

	cloud.Archive(cloud.ArchiveArgs{
		BodyText:          metadataJSONString,
		Key:               metadataS3Key,
		FileTypeExtension: "json",
	})

	bucket := GoTools.GetEnvVarHelper("S3_ARCHIVE_BUCKET")
	GoTools.Logger("INFO", fmt.Sprintf("Saved landing page metadata to s3://%s/%s", bucket, metadataS3Key))

	// TODO: filter stories?

	return LambdaResponse{
		OK:      true,
		Message: "Done",
	}, nil

	/*

		// e.g. 90 links in total
		// chunk size := 30
		// chunk count = 3
		chunkSize := 30
		chunkCount := int(math.Ceil(float64(len(stories) / chunkSize)))
		linkChunks := make([][]string, chunkCount)

		for i := 0; i < chunkCount; i++ {
			linkChunk := make([]string, chunkSize)
			for j := 0; j < chunkSize; j++ {
				linkChunk[j] = stories[i*chunkSize+j].URL
			}
			linkChunks[i] = linkChunk
		}

		GoTools.Logger("INFO", fmt.Sprintf("Pulled landing page content:\n ``` %s ``` \n ", landingPageHtmlText[:500]))

		for _, linkChunk := range linkChunks {
			// send SQS
			// refer to
			// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/sqs/sendmessage/
			queueName := GoTools.GetEnvVarHelper("STORIES_QUEUE_NAME")
			awsConfig, configErr := config.LoadDefaultConfig(context.TODO())
			if configErr != nil {
				GoTools.Logger("ERROR", "AWS shared configuration failed", configErr.Error())
			}
			sqsClient := sqs.NewFromConfig(awsConfig)
			getQueueResponse, getQueueError := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
			if getQueueError != nil {
				GoTools.Logger("ERROR", fmt.Sprintf("Error getting queue URL: %s", getQueueError.Error()))
			}
			queueURL := getQueueResponse.QueueUrl

			res, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
				// AWS required attributes
				// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_SendMessage.html
				// Golang API Reference
				// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs#SendMessageInput

				QueueUrl:    queueURL,
				MessageBody: aws.String(strings.Join(linkChunk, " ")),

				// Only FIFO queue can use `MessageGroupId`
				// MessageGroupId: aws.String(fmt.Sprintf("%s-00", queueName)),

				// TODO: add randomized delay
				DelaySeconds: 0,
			})

			if err != nil {
				GoTools.Logger("ERROR", fmt.Sprintf("Error sending message: %s", err))
			}
			GoTools.SimpleLogger("INFO", fmt.Sprintf("Message sent %s", *res.MessageId))
		}

		return LambdaResponse{
			OK:      true,
			Message: fmt.Sprintf("Sent %d messages OK", len(linkChunks)),
		}, nil

	*/
}

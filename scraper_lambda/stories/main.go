// You must use `main` package for lambda
// https://stackoverflow.com/a/50701572/9814131
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rivernews/GoTools"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)


func main() {
	lambda.Start(HandleRequest)
}

type LambdaEvent struct {
	LandingS3Key string `json:"landingS3Key"`
}

type LambdaResponse struct {
	OK bool `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, event LambdaEvent) (LambdaResponse, error) {
	GoTools.SendSlackMessage(fmt.Sprintf("Batch stories lambda started! Landing page S3 path: %s", event.LandingS3Key))

	// TODO: get all story links
	link := "http://story.com"

	// send SQS
	// refer to
	// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/sqs/sendmessage/
	queueName := GoTools.GetEnvVarHelper("STORIES_QUEUE_NAME")
	awsConfig, configErr := config.LoadDefaultConfig(context.TODO())
	if configErr != nil {
		GoTools.Logger("ERROR", "AWS shared configuration failed", configErr.Error())
	}
	sqsClient := sqs.NewFromConfig(awsConfig)
	getQueueResponse, getQueueError := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{ QueueName: aws.String(queueName) })
	if getQueueError != nil {
		GoTools.Logger("ERROR", fmt.Sprintf("Error getting queue URL: %s", getQueueError.Error()))
	}
	queueURL := getQueueResponse.QueueUrl

	res, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl: queueURL,
		MessageBody: aws.String(fmt.Sprintf("{\"storyURL\": \"%s\"}", link)),
		// TODO: better group id naming for multiplexing
		MessageGroupId: aws.String(fmt.Sprintf("%s-00", queueName)),
	})

	if err != nil {
		GoTools.Logger("ERROR", fmt.Sprintf("Error sending message: %s", err.Error()))
	}

	return LambdaResponse{
		OK: true,
		Message: fmt.Sprintf("Sent message OK: %s", *res.MessageId),
	}, nil
}

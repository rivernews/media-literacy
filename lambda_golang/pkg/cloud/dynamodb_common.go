package cloud

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rivernews/GoTools"

	"github.com/google/uuid"
)

var (
	dynamoDBClient     *dynamodb.Client
	dynamoDBClientOnce sync.Once
)

func SharedDynamoDBClient() *dynamodb.Client {
	dynamoDBClientOnce.Do(func() {
		dynamoDBClient = dynamodb.NewFromConfig(ShareAWSConfig())
	})
	return dynamoDBClient
}

func DynamoDBPutItem(ctx context.Context, item map[string]types.AttributeValue) *dynamodb.PutItemOutput {
	tableID := GoTools.GetEnvVarHelper("MEDIA_TABLE_ID")
	if tableID == "" {
		GoTools.Logger("ERROR", "MEDIA_TABLE_ID is required please set this env var")
	}

	if _, exist := item["uuid"]; !exist {
		item["uuid"] = uuid.New().String()
	}

	if _, exist := item["createdAt"]; !exist {
		// TODO
	}

	out, err := dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableID),
		Item:      item,
	})

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	return out
}

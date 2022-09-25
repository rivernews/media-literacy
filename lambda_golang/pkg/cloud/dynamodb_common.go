package cloud

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func DynamoDBPutItem(ctx context.Context, item map[string]any) *dynamodb.PutItemOutput {
	tableID := GoTools.GetEnvVarHelper("DYNAMODB_TABLE_ID")
	if tableID == "" {
		GoTools.Logger("ERROR", "DYNAMODB_TABLE_ID is required please set this env var")
	}

	if _, exist := item["uuid"]; !exist {
		item["uuid"] = uuid.New().String()
	}

	if _, exist := item["createdAt"]; !exist {
		// TODO
	}

	dynamoDBItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	out, err := dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableID),
		Item:      dynamoDBItem,
	})

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	return out
}

package cloud

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/newssite"

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

func getTableName() string {
	tableID := GoTools.GetEnvVarHelper("DYNAMODB_TABLE_ID")
	if tableID == "" {
		GoTools.Logger("ERROR", "DYNAMODB_TABLE_ID is required please set this env var")
	}
	return tableID
}

func DynamoDBPutItem(ctx context.Context, item any) *dynamodb.PutItemOutput {
	dynamoDBItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	if _, exist := dynamoDBItem["uuid"]; !exist {
		dynamoDBItem["uuid"] = uuid.New().String()
	}

	if _, exist := dynamoDBItem["createdAt"]; !exist {
		dynamoDBItem["createdAt"] = newssite.Now()
	}

	out, err := dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(getTableName()),
		Item:      dynamoDBItem,
	})

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	return out
}

func DynamoDBQueryWaitingMetadata(ctx context.Context, docType newssite.DocType) *[]newssite.MediaTableItem {
	out, err := dynamoDBClient.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(getTableName()),
		IndexName:              aws.String(string(newssite.METADATA_INDEX)),
		KeyConditionExpression: aws.String("isDocTypeWaitingForMetadata = :isDocTypeWaitingForMetadata and createdAt < :createdAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isDocTypeWaitingForMetadata": &types.AttributeValueMemberS{Value: string(docType)},
			":createdAt":                   &types.AttributeValueMemberS{Value: newssite.Now()},
		},
	})
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	var results []newssite.MediaTableItem
	attributevalue.UnmarshalListOfMaps(out.Items, &results)

	return &results
}

func DynamoDBUpdateItem() {
	// TODO: remove attribute

	// TODO: add event to `events`
}

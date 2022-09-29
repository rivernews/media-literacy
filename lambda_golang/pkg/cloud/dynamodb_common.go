package cloud

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/common"

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

func GetTableName() string {
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
		dynamoDBItem["uuid"] = &types.AttributeValueMemberS{Value: uuid.New().String()}
	}

	if _, exist := dynamoDBItem["createdAt"]; !exist {
		dynamoDBItem["createdAt"] = &types.AttributeValueMemberS{Value: common.Now()}
	}

	out, err := SharedDynamoDBClient().PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(GetTableName()),
		Item:      dynamoDBItem,
	})

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	return out
}

package newssite

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rivernews/GoTools"
	"github.com/rivernews/media-literacy/pkg/cloud"
	"github.com/rivernews/media-literacy/pkg/common"
)

type EventName string

const (
	// @`landing`  ✅ 📩 (put in db) ✅
	//EVENT_LANDING_PAGE_REQUESTED EventName = "LANDING_PAGE_REQUESTED"
	// -`landing_s3_trigger`
	EVENT_LANDING_PAGE_FETCHED EventName = "LANDING_PAGE_FETCHED"
	// @`landing_metadata` -> `landing_metadata_cronjob` (query db; store metadata) ✅ (cronjob trigger) ✅
	EVENT_LANDING_METADATA_REQUESTED EventName = "LANDING_METADATA_REQUESTED"
	EVENT_LANDING_METADATA_DONE      EventName = "LANDING_METADATA_DONE"
	// `stories` (metadata triggers; sfn) ✅
	EVENT_LANDING_STORIES_REQUESTED EventName = "LANDING_STORIES_REQUESTED"
	// `story` (sfn map; archive story) ✅ 📩
	//EVENT_STORY_REQUESTED EventName = "STORY_REQUESTED"
	// random wait
	EVENT_STORY_FETCHED EventName = "STORY_FETCHED"
	// +`stories_finalizer` (sfn last step)  ✅
	EVENT_LANDING_STORIES_FETCHED EventName = "LANDING_STORIES_FETCHED"
)

// func GetEventLandingPageRequested(newsSiteAlias string, newsSiteURL string) MediaTableItemEvent {
// 	return MediaTableItemEvent{
// 		EventName: EVENT_LANDING_PAGE_REQUESTED,
// 		Detail:    fmt.Sprintf("Requested landing page for %s at %s", newsSiteAlias, newsSiteURL),
// 		EventTime: common.Now(),
// 	}
// }

func GetEventLandingPageFetched(newsSiteAlias string, landingPageS3Key string) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_LANDING_PAGE_FETCHED,
		Detail:    fmt.Sprintf("Fetched landing page for %s; stored at %s", newsSiteAlias, landingPageS3Key),
		EventTime: common.Now(),
	}
}

func GetEventLandingMetadataRequested(landingPageS3Key string) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_LANDING_METADATA_REQUESTED,
		Detail:    fmt.Sprintf("Requested metadata for landing page %s", landingPageS3Key),
		EventTime: common.Now(),
	}
}

func GetEventLandingMetadataDone(metadataS3Key string, landingPageS3Key string) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_LANDING_METADATA_DONE,
		Detail:    fmt.Sprintf("Metadata is computed and archived at `%s`; for landing page `%s`", metadataS3Key, landingPageS3Key),
		EventTime: common.Now(),
	}
}

func GetEventLandingStoriesRequested(metadataS3Key string, storiesCount int) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_LANDING_STORIES_REQUESTED,
		Detail:    fmt.Sprintf("Stories(%d) requested for landing page based on metadata %s", storiesCount, metadataS3Key),
		EventTime: common.Now(),
	}
}

// if we want below, we need to PutItem w/o s3Key first
// then after fetched, UpdateItem - need to update s3Key as well!
// func GetEventStoryRequested(storyTitle string, storyURL string) MediaTableItemEvent {
// 	return MediaTableItemEvent{
// 		EventName: EVENT_STORY_REQUESTED,
// 		Detail:    fmt.Sprintf("Story %s requested %s", storyTitle, storyURL),
// 		EventTime: common.Now(),
// 	}
// }

func GetEventStoryFetched(storyTitle string, storyURL string) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_STORY_FETCHED,
		Detail:    fmt.Sprintf("Story %s fetched %s", storyTitle, storyURL),
		EventTime: common.Now(),
	}
}

func GetEventLandingStoriesFetched(landingPageS3Key string) MediaTableItemEvent {
	return MediaTableItemEvent{
		EventName: EVENT_LANDING_STORIES_FETCHED,
		Detail:    fmt.Sprintf("All stories fetched for landing page %s", landingPageS3Key),
		EventTime: common.Now(),
	}
}

type MediaTableItemEvent struct {
	EventName EventName `dynamodbav:"eventName" json:"eventName"`
	Detail    string    `dynamodbav:"detail" json:"detail"`
	EventTime string    `dynamodbav:"eventTime" json:"eventTime"`
}

type DocType string

const (
	DOCTYPE_LANDING DocType = "LANDING"
	DOCTYPE_STORY   DocType = "STORY"
)

type TableIndex string

const (
	METADATA_INDEX TableIndex = "metadataIndex"
	S3KEY_INDEX    TableIndex = "s3KeyIndex"
)

type MediaTableItem struct {
	Uuid                        string                `dynamodbav:"uuid,omitempty" json:"uuid,omitempty"`
	CreatedAt                   string                `dynamodbav:"createdAt,omitempty" json:"createdAt,omitempty"`
	S3Key                       string                `dynamodbav:"s3Key" json:"s3Key"`
	DocType                     DocType               `dynamodbav:"docType" json:"docType"`
	Events                      []MediaTableItemEvent `dynamodbav:"events" json:"events"`
	IsDocTypeWaitingForMetadata DocType               `dynamodbav:"isDocTypeWaitingForMetadata,omitempty" json:"isDocTypeWaitingForMetadata,omitempty"`
}

func DynamoDBQueryByS3Key(ctx context.Context, s3Key string) *[]MediaTableItem {
	out, err := cloud.SharedDynamoDBClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(cloud.GetTableName()),
		IndexName:              aws.String(string(S3KEY_INDEX)),
		KeyConditionExpression: aws.String("s3Key = :s3Key and createdAt < :createdAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":s3Key":     &types.AttributeValueMemberS{Value: s3Key},
			":createdAt": &types.AttributeValueMemberS{Value: common.Now()},
		},
		Limit: aws.Int32(10),
	})
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	var results []MediaTableItem
	attributevalue.UnmarshalListOfMaps(out.Items, &results)

	return &results
}

func DynamoDBQueryWaitingMetadata(ctx context.Context, docType DocType) *[]MediaTableItem {
	out, err := cloud.SharedDynamoDBClient().Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(cloud.GetTableName()),
		IndexName:              aws.String(string(METADATA_INDEX)),
		KeyConditionExpression: aws.String("isDocTypeWaitingForMetadata = :isDocTypeWaitingForMetadata and createdAt < :createdAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isDocTypeWaitingForMetadata": &types.AttributeValueMemberS{Value: string(docType)},
			":createdAt":                   &types.AttributeValueMemberS{Value: common.Now()},
		},
		Limit: aws.Int32(10),
	})
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	var results []MediaTableItem
	attributevalue.UnmarshalListOfMaps(out.Items, &results)

	return &results
}

func DynamoDBUpdateItem(ctx context.Context, uuid string, createdAt string, event MediaTableItemEvent, isMarkMetadataComplete bool) *dynamodb.UpdateItemOutput {
	dynamoDBItemEvent, err := attributevalue.MarshalMap(event)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}
	updateItemInput := dynamodb.UpdateItemInput{
		TableName: aws.String(cloud.GetTableName()),
		Key: map[string]types.AttributeValue{
			"uuid":      &types.AttributeValueMemberS{Value: uuid},
			"createdAt": &types.AttributeValueMemberS{Value: createdAt},
		},
		// manual
		// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html#Expressions.UpdateExpressions.ADD
		UpdateExpression: aws.String(`SET events = list_append(events, :e)`),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberM{Value: dynamoDBItemEvent},
				},
			},
		},
	}
	if isMarkMetadataComplete {
		*(updateItemInput.UpdateExpression) = *(updateItemInput.UpdateExpression) + ` REMOVE isDocTypeWaitingForMetadata`
	}

	out, err := cloud.SharedDynamoDBClient().UpdateItem(ctx, &updateItemInput)

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	return out
}

func DynamoDBUpdateItemAddEvent(ctx context.Context, uuid string, createdAt string, event MediaTableItemEvent) *dynamodb.UpdateItemOutput {
	return DynamoDBUpdateItem(ctx, uuid, createdAt, event, false)
}

func DynamoDBUpdateItemMarkAsMetadataComplete(ctx context.Context, uuid string, createdAt string, event MediaTableItemEvent) *dynamodb.UpdateItemOutput {
	return DynamoDBUpdateItem(ctx, uuid, createdAt, event, true)
}

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/rivernews/GoTools"

	"github.com/rivernews/media-literacy/pkg/cloud"
	"github.com/rivernews/media-literacy/pkg/common"
	"github.com/rivernews/media-literacy/pkg/newssite"
)

func main() {
	lambda.Start(HandleRequest)
}

type LambdaEvent struct {
	Name string `json:"name"`
}

type LambdaResponse struct {
	OK      bool   `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, name LambdaEvent) (LambdaResponse, error) {
	newsSite := newssite.GetNewsSite("NEWSSITE_ECONOMY")

	bodyText := common.Fetch(newsSite.LandingURL)
	GoTools.Logger("INFO", "In golang runtime now!\n\n```\n "+bodyText[:500]+"\n ...```\n End of message")

	// scraper
	result := newssite.GetStoriesFromEconomy(bodyText)

	// print out all story titles
	var slackMessage strings.Builder
	for i, topic := range result.Stories {
		slackMessage.WriteString(topic.Name)
		slackMessage.WriteString(" ")
		slackMessage.WriteString(topic.Description)
		slackMessage.WriteString(" ")
		slackMessage.WriteString(topic.URL)
		slackMessage.WriteString("\n")

		if i+1%50 == 0 {
			GoTools.SendSlackMessage(slackMessage.String())
			slackMessage.Reset()
		}
	}
	GoTools.SendSlackMessage(slackMessage.String())

	successMessage := fmt.Sprintf("Scraper finished - %d links found", len(result.Stories))
	GoTools.Logger("INFO", successMessage)

	// S3 archive
	landingPageS3Key := fmt.Sprintf("%s/daily-headlines/%s/landing.html", newsSite.Alias, common.Now())
	cloud.Archive(cloud.ArchiveArgs{
		BodyText: bodyText,
		Key:      landingPageS3Key,
	})

	return LambdaResponse{
		OK:      true,
		Message: "OK",
	}, nil
}

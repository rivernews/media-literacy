package main

import (
	"io"
	"fmt"
	"strings"
	"net/http"
	"golang.org/x/net/html/charset"
	"context"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/rivernews/GoTools"
)

func main() {
	lambda.Start(HandleRequest)
}

type LambdaEvent struct {
	Name string `json:"name"`
}

type LambdaResponse struct {
	OK bool `json:"OK:"`
	Message string `json:"message:"`
}

func HandleRequest(ctx context.Context, name LambdaEvent) (LambdaResponse, error) {
	newsSite := GetNewsSite("NEWSSITE_ECONOMY")
	resp, err := http.Get(newsSite.LandingURL)
	if err != nil {
		// handle error
		GoTools.Logger("ERROR", err.Error())
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type") // Optional, better guessing
	GoTools.Logger("INFO", "ContentType is ", contentType)
    utf8reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	body, err := io.ReadAll(utf8reader)
	if err != nil {
		// handle error
		GoTools.Logger("ERROR", err.Error())
	}
	bodyText := string(body)

	GoTools.Logger("INFO", "In golang runtime now!\n\n```\n " + bodyText[:500] + "\n ...```\n End of message")

	// scraper
	topics, err := getTopTenTrendingTopics(bodyText)

	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	var slackMessage strings.Builder
	for i, topic := range topics {
		slackMessage.WriteString(topic.Name)
		slackMessage.WriteString(" ")
		slackMessage.WriteString(topic.Description)
		slackMessage.WriteString(" ")
		slackMessage.WriteString(topic.URL)
		slackMessage.WriteString("\n")

		if i+1 % 50 == 0 {
			GoTools.SendSlackMessage(slackMessage.String())
			slackMessage.Reset()
		}
	}
	GoTools.SendSlackMessage(slackMessage.String())

	successMessage := fmt.Sprintf("Scraper finished - %d links found", len(topics))
	GoTools.Logger("INFO", successMessage)

	// S3 archive
	archive(strings.NewReader(bodyText), fmt.Sprintf("%s/daily-headlines/%s/landing.html", newsSite.Alias, time.Now().Format(time.RFC3339)))

	return LambdaResponse{
		OK: true,
		Message: "Slack command submitted successfully",
	}, nil
}

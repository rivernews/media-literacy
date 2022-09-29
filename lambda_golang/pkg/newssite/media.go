package newssite

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rivernews/GoTools"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type LandingPageMetadata struct {
	LandingPageS3Key     string  `json:"landingPageS3Key"`
	LandingPageUuid      string  `json:"landingPageUuid"`
	LandingPageCreatedAt string  `json:"landingPageCreatedAt"`
	Stories              []Topic `json:"stories"`
	UntitledStories      []Topic `json:"untitledstories"`
}

type StepFunctionInput struct {
	Stories              []Topic `json:"stories"`
	NewsSiteAlias        string  `json:"newsSiteAlias"`
	LandingPageUuid      string  `json:"landingPageUuid"`
	LandingPageS3Key     string  `json:"landingPageS3Key"`
	LandingPageTimeStamp string  `json:"landingPageTimeStamp"`
}

type StepFunctionMapIterationInput struct {
	Story                Topic  `json:"story"`
	NewsSiteAlias        string `json:"newsSiteAlias"`
	LandingPageUuid      string `json:"landingPageUuid"`
	LandingPageTimeStamp string `json:"landingPageTimeStamp"`
}

func GetStoriesFromEconomy(body string) *LandingPageMetadata {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	GoTools.Logger("INFO", "GoQuery "+string(doc.Find("title").Text()))

	// Find the review items
	var topics []Topic
	var untitledTopics []Topic
	// goquery api
	// https://pkg.go.dev/github.com/PuerkitoBio/goquery
	var emptyTitleURLs strings.Builder
	doc.Find("a[href$=html]").Each(func(i int, anchor *goquery.Selection) {
		topic := Topic{
			Name:        strings.ReplaceAll(strings.TrimSpace(anchor.Text()), "/", "-"),
			Description: "",
			URL:         strings.TrimSpace(anchor.AttrOr("href", "-")),
		}
		if topic.Name != "" {
			topics = append(topics, topic)
			GoTools.Logger("VERBOSE", "Found an anchor ", topic.Name, topic.URL)
		} else {
			untitledTopics = append(untitledTopics, topic)
			emptyTitleURLs.WriteString(topic.URL)
			emptyTitleURLs.WriteString("\n")
		}
	})
	GoTools.Logger("INFO", "Skipped due to empty title URLs:\n", emptyTitleURLs.String())

	return &LandingPageMetadata{
		Stories:         topics,
		UntitledStories: untitledTopics,
	}
}

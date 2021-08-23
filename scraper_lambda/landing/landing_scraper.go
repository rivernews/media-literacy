package main

import (
	"strings"

	"github.com/rivernews/GoTools"
	"github.com/PuerkitoBio/goquery"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func getTopTenTrendingTopics(body string) ([]Topic, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	GoTools.Logger("INFO", "GoQuery " + string(doc.Find("title").Text()))

	// Find the review items
	var topics []Topic
	// goquery api
	// https://pkg.go.dev/github.com/PuerkitoBio/goquery
	var emptyTitleURLs strings.Builder
	doc.Find("a[href$=html]").Each(func(i int, anchor *goquery.Selection) {
		topic := Topic{
			Name: anchor.Text(),
			Description: "",
			URL: anchor.AttrOr("href", "-"),
		}
		if topic.Name != "" {
			topics = append(topics, topic)
			GoTools.Logger("DEBUG", "Found an anchor ", topic.Name, topic.URL)
		} else {
			emptyTitleURLs.WriteString(topic.URL)
			emptyTitleURLs.WriteString("\n")
		}
	})
	GoTools.Logger("INFO", "Skipped due to empty title URLs:\n", emptyTitleURLs.String())

	return topics, nil
}

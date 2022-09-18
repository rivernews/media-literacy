package newssite

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"

	"github.com/rivernews/GoTools"
)

type NewsSite struct {
	Key        string
	Name       string
	Alias      string
	LandingURL string
}

func GetNewsSite(envVar string) NewsSite {
	ssmValue := GoTools.GetEnvVarHelper(envVar)

	tokens := strings.Split(ssmValue, ",")

	return NewsSite{
		Key:        tokens[0],
		Name:       tokens[1],
		Alias:      tokens[2],
		LandingURL: tokens[3],
	}
}

func Fetch(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		GoTools.Logger("ERROR", err.Error())
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type") // Optional, better guessing
	GoTools.Logger("DEBUG", "ContentType is ", contentType)
	utf8reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}

	body, err := io.ReadAll(utf8reader)
	if err != nil {
		// handle error
		GoTools.Logger("ERROR", err.Error())
	}
	return string(body)
}

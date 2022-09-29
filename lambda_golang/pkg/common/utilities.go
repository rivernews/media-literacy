package common

import (
	"io"
	"net/http"
	"time"

	"github.com/rivernews/GoTools"
	"golang.org/x/net/html/charset"
)

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

func Now() string {
	return time.Now().Format(time.RFC3339)
}

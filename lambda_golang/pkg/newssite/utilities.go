package newssite

import (
	"encoding/json"
	"strings"

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

func AsJson(v any) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}
	return string(jsonBytes)
}

func FromJson(b []byte, structInstance any) {
	if err := json.Unmarshal(b, structInstance); err != nil {
		GoTools.Logger("ERROR", err.Error())
	}
}

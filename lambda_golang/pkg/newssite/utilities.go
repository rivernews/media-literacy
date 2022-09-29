package newssite

import (
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

package cloud

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/rivernews/GoTools"

	"context"
	"sync"
)

var (
	awsConfig      aws.Config
	awsConfigOnce  sync.Once
	awsConfigError error
)

func ShareAWSConfig() aws.Config {
	awsConfigOnce.Do(func() {
		// Based on
		// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
		awsConfig, awsConfigError = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion("us-west-2"),
		)
		if awsConfigError != nil {
			GoTools.Logger("ERROR", "AWS shared configuration failed", awsConfigError.Error())
		}
	})
	return awsConfig
}

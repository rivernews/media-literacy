package cloud

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/rivernews/GoTools"

	"context"
	"sync"
)

var (
	s3Client *s3.Client
	s3ClientOnce sync.Once
)

func SharedS3Client() *s3.Client {
	// followed
	// https://stackoverflow.com/a/53504651/9814131
	s3ClientOnce.Do(func() {
		// Based on
		// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
		awsConfig, configErr := config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion("us-west-2"),
		)
		if configErr != nil {
			GoTools.Logger("ERROR", "AWS shared configuration failed", configErr.Error())
		}

		s3Client = s3.NewFromConfig(awsConfig)
	})

	return s3Client
}
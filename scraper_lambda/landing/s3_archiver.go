package main

import (
	"context"
  	"io"
  	"fmt"
  	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/rivernews/GoTools"
)

func archive(body io.Reader, key string) (bool, error) {
	bucket := GoTools.GetEnvVarHelper("S3_ARCHIVE_BUCKET")
	GoTools.Logger("INFO", "Bucket to archive: s3://", bucket, "Key:", key)

	// Based on
	// https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
	awsConfig, configErr := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if configErr != nil {
		GoTools.Logger("ERROR", "AWS shared configuration failed", configErr.Error())
	}

  	timeout := time.Second * 30

  	client := s3.NewFromConfig(awsConfig)

	// Based on
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#Uploader
	uploader := manager.NewUploader(client)

  	// Create a context with a timeout that will abort the upload if it takes
  	// more than the passed in timeout.
  	ctx := context.Background()
  	var cancelFn func()
  	if timeout > 0 {
  		ctx, cancelFn = context.WithTimeout(ctx, timeout)
  	}
  	// Ensure the context is canceled to prevent leaking.
  	// See context package for more information, https://golang.org/pkg/context/
	if cancelFn != nil {
  		defer cancelFn()
	}

  	// Uploads the object to S3. The Context will interrupt the request if the
  	// timeout expires.
  	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
  		Bucket: aws.String(bucket),
  		Key:    aws.String(key),
  		Body:   body,
		ContentType: aws.String("text/html"),
  	})
  	if err != nil {
		GoTools.Logger("ERROR", fmt.Sprintf("failed to upload object: %v", err))
  	}

  	GoTools.Logger("INFO", fmt.Sprintf("successfully uploaded file to `s3://%s/%s`\n", bucket, key))

	return true, nil
}

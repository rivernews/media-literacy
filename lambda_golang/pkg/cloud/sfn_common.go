package cloud

import (
	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/rivernews/GoTools"

	"context"
	"sync"
)

var (
	sfnClient     *sfn.Client
	sfnClientOnce sync.Once
)

func SharedSfnClient() *sfn.Client {
	sfnClientOnce.Do(func() {
		sfnClient = sfn.NewFromConfig(ShareAWSConfig())
	})
	return sfnClient
}

// refer to real example https://qiita.com/tanaka_takurou/items/bba2e6c32052a636a272
// SDK doc https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sfn#Client.StartExecution
func SfnStartExecution(ctx context.Context, sfnInput *sfn.StartExecutionInput) *sfn.StartExecutionOutput {
	sfnClient = SharedSfnClient()

	output, err := sfnClient.StartExecution(ctx, sfnInput)
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}
	return output
}

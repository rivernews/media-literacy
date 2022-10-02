package cloud

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"sync"
)

var (
	s3Client     *s3.Client
	s3ClientOnce sync.Once
)

func SharedS3Client() *s3.Client {
	// followed
	// https://stackoverflow.com/a/53504651/9814131
	s3ClientOnce.Do(func() {
		s3Client = s3.NewFromConfig(ShareAWSConfig())
	})

	return s3Client
}

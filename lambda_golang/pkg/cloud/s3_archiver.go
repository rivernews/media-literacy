package cloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/rivernews/GoTools"
)

type ArchiveArgs struct {
	BodyText          string
	Key               string
	FileTypeExtension string
}

func Archive(args ArchiveArgs) (bool, error) {
	fileTypeExtension := "text/html"
	if strings.ToLower(args.FileTypeExtension) == "json" {
		fileTypeExtension = "application/json"
	}

	bucket := GoTools.GetEnvVarHelper("S3_ARCHIVE_BUCKET")
	GoTools.Logger("INFO", "Bucket to archive: `s3://", bucket, "` Key: `", args.Key, "`")

	timeout := time.Second * 30
	client := SharedS3Client()

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
		Bucket:      aws.String(bucket),
		Key:         aws.String(args.Key),
		Body:        strings.NewReader(args.BodyText),
		ContentType: aws.String(fileTypeExtension),
	})
	if err != nil {
		GoTools.Logger("ERROR", fmt.Sprintf("failed to upload object: %v", err))
	}

	GoTools.Logger("INFO", fmt.Sprintf("successfully uploaded file to `s3://%s/%s`\n", bucket, args.Key))

	return true, nil
}

func IsDuplicated(key string) bool {
	exist, _ := Exist(key)

	return exist
}

func Exist(key string) (bool, *s3.HeadObjectOutput) {
	bucket := GoTools.GetEnvVarHelper("S3_ARCHIVE_BUCKET")
	client := SharedS3Client()

	// based on
	// https://stackoverflow.com/a/65710928/9814131
	headObject, headError := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if headError != nil {
		headErrorMessage := headError.Error()
		if strings.Contains(headErrorMessage, "404") {
			return false, nil
		}

		GoTools.Logger("ERROR", headError.Error())
	}

	return true, headObject
}

func Pull(key string) string {
	bucket := GoTools.GetEnvVarHelper("S3_ARCHIVE_BUCKET")
	client := SharedS3Client()

	// based on
	// https://stackoverflow.com/a/65710928/9814131
	_, headObject := Exist(key)

	// main idea
	// https://stackoverflow.com/a/41645765/9814131
	// code based on
	// https://github.com/aws/aws-sdk-go-v2/pull/1171/files#diff-c43ccf2f39bfbd136d7f7ddf2a1c88ac983d910b687bca29b4a8e6ea9759551b
	// pre-allocate in memory buffer, where headObject type is *s3.HeadObjectOutput
	// and
	// AWS SDK v2 Doc
	// https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/#download-manager

	downloader := manager.NewDownloader(client)
	buf := make([]byte, int(headObject.ContentLength))
	// wrap with aws.WriteAtBuffer
	w := manager.NewWriteAtBuffer(buf)
	// download file into the memory
	numBytesDownloaded, err := downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		GoTools.Logger("ERROR", err.Error())
	}
	GoTools.Logger("INFO", fmt.Sprintf("Downloaded %d for `s3://%s/%s`", numBytesDownloaded, bucket, key))

	return string(w.Bytes())
}

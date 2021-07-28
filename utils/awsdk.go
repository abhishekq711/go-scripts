package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"go.uber.org/zap"
)

func DownloadObject(bucket, item, region string) error {

	file, err := os.Create(item)
	if err != nil {
		return err
	}

	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		os.Remove(item)
		return err
	}
	zap.L().Info(fmt.Sprintf("Downloaded %s %d bytes", file.Name(), numBytes))
	return nil
}

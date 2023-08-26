package destination

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type S3DestinationParams struct {
	Name   string
	Bucket string
	Region string
}

type s3BucketDestination struct {
	params S3DestinationParams
}

func (s *s3BucketDestination) String() string {
	return fmt.Sprintf("Name: %s, Type: %s", s.params.Name, "aws")
}

func NewS3BucketDestination(params S3DestinationParams) Destination {
	return &s3BucketDestination{
		params: params,
	}
}

func (s *s3BucketDestination) UploadFiles(fileList []string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s.params.Region))
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg)

	prefix := strings.SplitN(s.params.Bucket, "/", 2)

	for _, f := range fileList {
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		slog.Info("Upload file", "file", filepath.Base(f))
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &s.params.Bucket,
			Key:    aws.String(fmt.Sprintf("%s/%s", prefix[1], filepath.Base(file.Name()))),
			Body:   file,
		})

		if err != nil {
			return err
		}
		file.Close()
	}

	if err := s.cleanupLocalBackup(fileList); err != nil {
		return err
	}

	return nil
}

func (s *s3BucketDestination) cleanupLocalBackup(fileList []string) error {
	for _, f := range fileList {
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	return nil
}

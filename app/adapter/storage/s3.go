package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"pixelPromo/config"
	"pixelPromo/domain/port"
)

type bucketS3Storage struct {
	api *s3.Client
	cfg *config.Config
}

func NewBucketS3Storage(
	cfgAws *aws.Config,
	cfg *config.Config,
) port.Storage {

	api := s3.NewFromConfig(*cfgAws, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	return &bucketS3Storage{
		api: api,
		cfg: cfg,
	}
}

func (b bucketS3Storage) UploadUserPicture(ctx context.Context, fileName string, reader io.Reader) (string, error) {
	uploader := manager.NewUploader(b.api)

	bucketName := b.cfg.Viper.GetString("aws.s3.buckets.user-pictures")
	out, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	})

	if err != nil {
		return "", err
	}

	if out == nil {
		return "", fmt.Errorf("failed to upload file")
	}

	return out.Location, nil
}

func (b bucketS3Storage) UploadPromotionImage(ctx context.Context, fileName string, reader io.Reader) (string, error) {
	uploader := manager.NewUploader(b.api)

	bucketName := b.cfg.Viper.GetString("aws.s3.buckets.promotion-images")
	out, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	})

	if err != nil {
		return "", err
	}

	if out == nil {
		return "", fmt.Errorf("failed to upload file")
	}

	return out.Location, nil
}

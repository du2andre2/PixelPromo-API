package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"io/ioutil"
)

type BucketS3 interface {
	GetFile(ctx context.Context, input *GetFileInput) (*GetFileOutput, error)
	PutFile(ctx context.Context, input *PutFileInput) (*PutFileOutput, error)
}

type bucketS3 struct {
	api *s3.Client
}

func NewBucketS3(cfg *aws.Config) BucketS3 {

	api := s3.NewFromConfig(*cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	return &bucketS3{
		api: api,
	}
}

type GetFileInput struct {
	BucketName string
	FileName   string
}

type GetFileOutput struct {
	BodyItem []byte
}

func (d *bucketS3) GetFile(ctx context.Context, input *GetFileInput) (*GetFileOutput, error) {
	resp, err := d.api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(input.BucketName),
		Key:    aws.String(input.FileName),
	})
	if err != nil {
		return nil, err
	}

	s3objectBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &GetFileOutput{
		BodyItem: s3objectBytes,
	}, nil
}

type PutFileInput struct {
	BucketName string
	FileName   string
	BodyItem   io.Reader
}

type PutFileOutput struct {
	UrlFile string
}

func (d *bucketS3) PutFile(ctx context.Context, input *PutFileInput) (*PutFileOutput, error) {
	_, err := d.api.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(input.BucketName),
		Key:    aws.String(input.FileName),
		Body:   input.BodyItem,
	})
	if err != nil {
		return nil, err
	}

	return &PutFileOutput{}, nil
}

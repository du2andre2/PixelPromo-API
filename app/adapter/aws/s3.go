package aws

import "github.com/aws/aws-sdk-go-v2/aws"

type BucketS3 interface {
	GetFile() *File
}

type bucketS3 struct {
	cfg *aws.Config
}

func NewBucketS3(cfg *aws.Config) BucketS3 {
	return &bucketS3{
		cfg: cfg,
	}
}

type File struct{}

func (d *bucketS3) GetFile() *File {
	return &File{}
}

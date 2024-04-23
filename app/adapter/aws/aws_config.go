package aws

import "github.com/aws/aws-sdk-go-v2/aws"

func NewConfigAWS() *aws.Config {
	return awsConfig()
}

func awsConfig() *aws.Config {
	return &aws.Config{}
}

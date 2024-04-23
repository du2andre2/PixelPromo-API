package aws

import "github.com/aws/aws-sdk-go-v2/aws"

type DynamoDB interface {
	GetItem() *Item
}

type dynamoDB struct {
	cfg *aws.Config
}

func NewDynamoDb(cfg *aws.Config) DynamoDB {
	return &dynamoDB{
		cfg: cfg,
	}
}

type Item struct{}

func (d *dynamoDB) GetItem() *Item {
	return &Item{}
}

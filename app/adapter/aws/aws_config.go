package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"log"
	cfg "pixelPromo/adapter/config"
)

func NewConfigAWS(c *cfg.Config) *aws.Config {

	if c.Env == cfg.Local {

		awsConfig, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Viper.GetString("aws.config.region")),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("local", "local", "")),
		)

		endpoint := c.Viper.GetString("aws.config.local-endpoint")
		awsConfig.BaseEndpoint = &endpoint

		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		return &awsConfig
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Viper.GetString("aws.config.region")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &awsConfig
}

package instance

import (
	stackConfig "bitbucket.org/codapayments/coda-stack-management-app/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"context"
)

func GetConnection() *dynamodb.Client {

	configs := stackConfig.GetConfig()
	var region string
	var endpoint string

	if configs.Environment != "" && configs.Environment == "development" {
		region = "localhost"
		endpoint = "http://localhost:8000"

	}else{
		region = configs.Region
		endpoint = "https://dynamodb." + region + ".amazonaws.com"
	}
	
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))

	if err != nil {
		panic(err)
	}

	// Create DynamoDB client 
	svc := dynamodb.NewFromConfig(cfg)

	// Create DynamoDB client
	return svc
}

package providers

import (
	"autovpn/options"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type AWS struct {
	Provider
}

type AWSOptions struct {
}

func getClient(avpnConfig options.Config) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", "")))
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(cfg), nil
}

func (A AWS) GetRegions(args ProviderArgs) ([]Region, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) GetInstances(args ProviderArgs) ([]Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) CreateServer(args ProviderArgs) (*Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) AwaitProvisioning(args ProviderArgs) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) DestroyServer(args ProviderArgs) error {
	//TODO implement me
	panic("implement me")
}

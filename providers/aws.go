package providers

import (
	"autovpn/data"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type AWS struct {
	Client
}

type AWSOptions struct {
}

func getClient(avpnConfig data.Config) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", "")))
	if err != nil {
		return nil, err
	}
	return ec2.NewFromConfig(cfg), nil
}

func (A AWS) getRegions(args data.ArgsBundle) ([]Region, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) getInstances(args data.ArgsBundle) ([]data.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) createServer(args data.ArgsBundle) (*data.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) awaitProvisioning(args data.ArgsBundle) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) destroyServer(args data.ArgsBundle) error {
	//TODO implement me
	panic("implement me")
}

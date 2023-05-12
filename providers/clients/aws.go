package clients

import (
	"autovpn/data"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var _ Client = &AWS{}

type AWS struct{}

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

func (A AWS) GetRegions(args data.ArgsBundle) ([]Region, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) GetInstances(args data.ArgsBundle) ([]data.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) CreateServer(args data.ArgsBundle) (*data.Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) AwaitProvisioning(args data.ArgsBundle) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) DestroyServer(args data.ArgsBundle) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) Connect(args data.ArgsBundle) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) TimeoutSetup(args data.ArgsBundle) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

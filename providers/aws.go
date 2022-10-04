package providers

import "autovpn/options"

type AWS struct {
	Provider
}

func (A AWS) GetRegions() ([]Region, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) GetInstances(config options.Config) ([]Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) CreateServer(arguments options.Arguments, config options.Config) (*Instance, error) {
	//TODO implement me
	panic("implement me")
}

func (A AWS) AwaitProvisioning(instance Instance, token string) error {
	//TODO implement me
	panic("implement me")
}

func (A AWS) DestroyServer(instance Instance, token string) error {
	//TODO implement me
	panic("implement me")
}

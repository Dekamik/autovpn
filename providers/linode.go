package providers

type Linode struct {
	Provider
}

func (l Linode) GetRegions() ([]Region, error) {
	// TODO: Implement
	return nil, nil
}

func (l Linode) CreateServer() (Instance, error) {
	// TODO: Implement
	return Instance{}, nil
}

func (l Linode) DestroyServer(instance Instance) error {
	// TODO: Implement
	return nil
}

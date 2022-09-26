package providers

import (
	"errors"
	"fmt"
)

var AvailableProviders = []string{
	"linode",
}

type Provider interface {
	GetRegions() ([]Region, error)
	CreateServer() (Instance, error)
	DestroyServer(Instance) error
}

type Region struct {
	Id      string
	Country string
}

type Instance struct {
	Id string
}

func New(name string) (Provider, error) {
	switch name {

	case "linode":
		return Linode{}, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unknown provider %s", name))
	}
}

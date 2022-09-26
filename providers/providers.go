package providers

import (
	"autovpn/options"
	"fmt"
)

var AvailableProviders = []string{
	"linode",
}

type Provider interface {
	GetRegions(silent bool) ([]Region, error)
	CreateServer(arguments options.Arguments, config options.Config) (*Instance, error)
	DestroyServer(instance Instance, authHeader string) error
}

type Region struct {
	Id      string
	Country string
}

type Instance struct {
	Id        string
	IpAddress string
}

func New(name string) (Provider, error) {
	switch name {

	case "linode":
		return Linode{}, nil

	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}

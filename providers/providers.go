package providers

import (
	"autovpn/config"
	"fmt"
)

var AvailableProviders = []string{
	"linode",
}

type Provider interface {
	GetRegions(silent bool) ([]Region, error)
	CreateServer(arguments config.Arguments, providerConfig ProviderConfig) (*Instance, error)
	DestroyServer(instance Instance, authHeader string) error
}

type ProviderConfig interface {
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

func GetProviderConfig(yamlConfig config.YamlConfig, name string) (ProviderConfig, error) {
	switch name {

	case "linode":
		providerConfig := yamlConfig.Providers[name].(map[string]interface{})
		return LinodeConfig{
			Image:    providerConfig["image"].(string),
			Key:      providerConfig["key"].(string),
			TypeSlug: providerConfig["type_slug"].(string),
		}, nil

	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}

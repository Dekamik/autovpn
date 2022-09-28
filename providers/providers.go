package providers

import (
	"autovpn/options"
	"fmt"
)

// Add implemented providers here.
var availableProviders = map[string]Provider{
	"linode": Linode{},
}

type Region struct {
	Id      string
	Country string
}

type Instance struct {
	Id        string
	IpAddress string
	RootUser  string
	RootPass  string
	SshPort   int
}

type Provider interface {
	// GetRegions downloads available server regions for the provider.
	GetRegions() ([]Region, error)

	// CreateServer creates, provisions and boots the server in the cloud.
	CreateServer(arguments options.Arguments, config options.Config) (*Instance, error)

	// AwaitProvisioning blocks the thread until the server is ready to receive SSH connections.
	AwaitProvisioning(instance Instance, token string) error

	// DestroyServer destroys the server.
	DestroyServer(instance Instance, token string) error
}

func New(name string) (Provider, error) {
	provider := availableProviders[name]

	if provider == nil {
		return nil, fmt.Errorf("unknown provider %q", name)
	}

	return provider, nil
}

func ListProviders() []string {
	list := make([]string, len(availableProviders))
	i := 0
	for name := range availableProviders {
		list[i] = name
		i++
	}
	return list
}

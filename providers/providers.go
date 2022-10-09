package providers

import (
	"autovpn/options"
	"fmt"
)

const (
	InstanceTag = "AutoVPN Instance"
)

// Add implemented providers here.
var availableProviders = map[string]Provider{
	"aws":    AWS{},
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
	Tags      []string
}

type ProviderArgs struct {
	Config    options.Config
	Arguments options.Arguments
	Instance  Instance
}

type Provider interface {
	// GetRegions downloads available server regions for the provider.
	GetRegions(args ProviderArgs) ([]Region, error)

	// GetInstances downloads all AutoVPN instances at the provider.
	GetInstances(args ProviderArgs) ([]Instance, error)

	// CreateServer creates, provisions and boots the server in the cloud.
	CreateServer(args ProviderArgs) (*Instance, error)

	// AwaitProvisioning blocks the thread until the server is ready to receive SSH connections.
	AwaitProvisioning(args ProviderArgs) error

	// DestroyServer destroys the server.
	DestroyServer(args ProviderArgs) error
}

func NewProvider(name string) (Provider, error) {
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

package clients

import (
	"autovpn/data"
)

const (
	InstanceTag = "AutoVPN Instance"
)

// AvailableProviders contains all available provider names and their corresponding clients.
// Add implemented providers here.
var AvailableProviders = map[string]Client{
	//"aws":    AWS{}, TODO: Uncomment
	"linode": Linode{},
}

type Region struct {
	Id      string
	Country string
}

type Client interface {
	// GetRegions downloads available server regions for the provider.
	GetRegions(args data.ArgsBundle) ([]Region, error)

	// GetInstances downloads all AutoVPN instances at the provider.
	GetInstances(args data.ArgsBundle) ([]data.Instance, error)

	// CreateServer creates, provisions and boots the server in the cloud.
	CreateServer(args data.ArgsBundle) (*data.Instance, error)

	// AwaitProvisioning blocks the thread until the server is ready to receive SSH connections.
	AwaitProvisioning(args data.ArgsBundle) error

	// DestroyServer destroys the server.
	DestroyServer(args data.ArgsBundle) error

	// Connect is a helper method that provisions and connects to a provider's client, if required (example: AWS Client)
	Connect(args data.ArgsBundle) error

	// TimeoutSetup returns the required bash commands to set up automatic deletion on server
	TimeoutSetup(args data.ArgsBundle) ([]string, error)
}

func New(providerName string, args data.ArgsBundle) (*Client, error) {
	client := AvailableProviders[providerName]

	err := client.Connect(args)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

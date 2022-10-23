package providers

import (
	"autovpn/data"
)

const (
	InstanceTag = "AutoVPN Instance"
)

// Add implemented providers here.
var availableProviders = map[string]Client{
	//"aws":    AWS{}, TODO: Uncomment
	"linode": Linode{},
}

type Region struct {
	Id      string
	Country string
}

type Client interface {
	// getRegions downloads available server regions for the provider.
	getRegions(args data.ArgsBundle) ([]Region, error)

	// getInstances downloads all AutoVPN instances at the provider.
	getInstances(args data.ArgsBundle) ([]data.Instance, error)

	// createServer creates, provisions and boots the server in the cloud.
	createServer(args data.ArgsBundle) (*data.Instance, error)

	// awaitProvisioning blocks the thread until the server is ready to receive SSH connections.
	awaitProvisioning(args data.ArgsBundle) error

	// destroyServer destroys the server.
	destroyServer(args data.ArgsBundle) error

	// connect is a helper method that provisions and connects to a provider's client, if required (example: AWS Client)
	connect(args data.ArgsBundle) error

	// failSafeSetup returns the required bash commands to set up automatic deletion on server
	failSafeSetup(args data.ArgsBundle) ([]string, error)
}

func newClient(providerName string, args data.ArgsBundle) (*Client, error) {
	client := availableProviders[providerName]

	err := client.connect(args)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

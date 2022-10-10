package providers

import (
	"autovpn/data"
	"autovpn/options"
)

const (
	InstanceTag = "AutoVPN Instance"
)

// Add implemented providers here.
var availableProviders = map[string]Client{
	//"aws":    AWS{},
	"linode": Linode{},
}

type Region struct {
	Id      string
	Country string
}

type ClientArgs struct {
	Config    options.Config
	Arguments options.Arguments
	Instance  data.Instance
}

type Client interface {
	// getRegions downloads available server regions for the provider.
	getRegions(args ClientArgs) ([]Region, error)

	// getInstances downloads all AutoVPN instances at the provider.
	getInstances(args ClientArgs) ([]data.Instance, error)

	// createServer creates, provisions and boots the server in the cloud.
	createServer(args ClientArgs) (*data.Instance, error)

	// awaitProvisioning blocks the thread until the server is ready to receive SSH connections.
	awaitProvisioning(args ClientArgs) error

	// destroyServer destroys the server.
	destroyServer(args ClientArgs) error

	// connect
	connect(args ClientArgs) error
}

func newClient(providerName string, args ClientArgs) (*Client, error) {
	client := availableProviders[providerName]

	err := client.connect(args)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

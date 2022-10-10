package providers

import (
	"autovpn/data"
	"autovpn/helpers"
	"autovpn/openvpn"
	"fmt"
	"log"
	"os"
	"sync"
)

type ProviderInterface interface {
	// Connect provisions a new server and connects to it.
	Connect() error

	ListZombies() error

	Purge() error

	// ShowRegions downloads available server regions for the provider.
	ShowRegions() error
}

type Provider struct {
	client Client
	args   ClientArgs
}

func destroyServer(client Client, args ClientArgs) {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Destroying server", finish, exited)
	err := client.destroyServer(args)
	finish <- true
	<-exited
	if err != nil {
		panic(err)
	}
}

func removeOvpnConfig(ovpnConfig *string) {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint(fmt.Sprintf("Removing %s", *ovpnConfig), finish, exited)
	_ = os.Remove(*ovpnConfig)
	finish <- true
	<-exited
}

func (p Provider) Connect() error {
	if !p.args.Arguments.NoAdminCheck {
		isAdmin, err := helpers.IsAdmin()
		if err != nil {
			return err
		}
		if !isAdmin {
			return fmt.Errorf("Root/Administrator privileges required")
		}
	}

	exe := openvpn.GetExecutable(p.args.Config.Overrides.OpenvpnExe)
	if isInstalled := openvpn.IsInstalled(exe); !isInstalled {
		return fmt.Errorf("couldn't find OpenVPN exe (%s). OpenVPN must be installed", exe)
	}

	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Creating instance", finish, exited)
	instance, err := p.client.createServer(p.args)
	p.args.Instance = *instance
	finish <- true
	<-exited
	if err != nil {
		if instance != nil {
			destroyServer(p.client, p.args)
		}
		return err
	}
	defer destroyServer(p.client, p.args)

	go helpers.WaitPrint("Starting instance", finish, exited)
	err = p.client.awaitProvisioning(p.args)
	finish <- true
	<-exited
	if err != nil {
		return err
	}

	go helpers.WaitPrint("Installing OpenVPN Server", finish, exited)
	ovpnConfig, err := openvpn.Install(*instance, p.args.Config.Agent.ScriptUrl)
	finish <- true
	<-exited
	if err != nil {
		return err
	}
	defer removeOvpnConfig(ovpnConfig)

	err = openvpn.Connect(exe, *ovpnConfig)
	if err != nil {
		return err
	}

	return nil
}

func (p Provider) ListZombies() error {
	instances, err := p.client.getInstances(p.args)
	providerName := p.args.Arguments.Provider
	if err != nil {
		return err
	}

	fmt.Printf("--- %s: %d ---\n", providerName, len(instances))
	for _, instance := range instances {
		fmt.Printf("%s %s\n", providerName, instance.Id)
	}

	return nil
}

func (p Provider) Purge() error {
	instances, err := p.client.getInstances(p.args)
	providerName := p.args.Arguments.Provider
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, instance := range instances {
		isAutoVPNInstance := false
		for _, tag := range instance.Tags {
			if tag == InstanceTag {
				isAutoVPNInstance = true
				break
			}
		}

		if !isAutoVPNInstance {
			continue
		}

		wg.Add(1)
		go func(instance data.Instance) {
			log.Printf("Purging %s %s...", providerName, instance.Id)
			err := p.client.destroyServer(p.args)
			if err != nil {
				log.Fatalf("Purge ERR %s %s: %s", providerName, instance.Id, err.Error())
			}
			log.Printf("Purge OK %s %s", providerName, instance.Id)
			wg.Done()
		}(instance)
	}

	wg.Wait()
	return nil
}

func (p Provider) ShowRegions() error {
	regions, err := p.client.getRegions(p.args)
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Printf("%s (%s)\n", region.Id, region.Country)
	}

	return nil
}

func New(providerName string, args ClientArgs) (*Provider, error) {
	client, err := newClient(providerName, args)
	if err != nil {
		return nil, err
	}

	return &Provider{
		client: *client,
		args:   args,
	}, nil
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

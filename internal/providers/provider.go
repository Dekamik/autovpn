package providers

import (
	"autovpn/internal/data"
	"autovpn/internal/helpers"
	"autovpn/internal/openvpn"
	"autovpn/internal/providers/clients"
	"fmt"
	"log"
	"os"
	"sync"
)

type ProviderInterface interface {
	// Connect provisions a new VPN server and connects to it.
	Connect() error

	// ListZombies lists all existing AutoVPN servers on the provider.
	ListZombies() error

	// Purge destroys all AutoVPN servers on the provider.
	Purge() error

	// ShowRegions downloads available server regions for the provider.
	ShowRegions() error
}

type Provider struct {
	name   string
	client clients.Client
	args   data.ArgsBundle
}

func destroyServer(client clients.Client, args data.ArgsBundle) {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Destroying server", finish, exited)
	err := client.DestroyServer(args)
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
	instance, err := p.client.CreateServer(p.args)
	if err != nil {
		return err
	}
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
	err = p.client.AwaitProvisioning(p.args)
	finish <- true
	<-exited
	if err != nil {
		return err
	}

	go helpers.WaitPrint("Installing OpenVPN Server", finish, exited)
	timeoutSetup, err := p.client.TimeoutSetup(p.args)
	ovpnConfig, err := openvpn.Install(*instance, p.args.Config.Agent.ScriptUrl, timeoutSetup)
	finish <- true
	<-exited
	if err != nil {
		return err
	}
	defer removeOvpnConfig(ovpnConfig)

	err = openvpn.Connect(exe, *ovpnConfig, timeoutSetup, *instance)
	if err != nil {
		return err
	}

	return nil
}

func (p Provider) ListZombies() error {
	instances, err := p.client.GetInstances(p.args)
	if err != nil {
		return err
	}

	fmt.Printf("--- %s: %d ---\n", p.name, len(instances))
	for _, instance := range instances {
		fmt.Printf("%s %s\n", p.name, instance.Id)
	}

	return nil
}

func (p Provider) Purge() error {
	instances, err := p.client.GetInstances(p.args)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, instance := range instances {
		for _, tag := range instance.Tags {
			if tag == clients.InstanceTag {
				wg.Add(1)
				go func(instance data.Instance) {
					log.Printf("Purging %s %s...", p.name, instance.Id)
					err := p.client.DestroyServer(p.args)
					if err != nil {
						log.Fatalf("Purge ERR %s %s: %s", p.name, instance.Id, err.Error())
					}
					log.Printf("Purge OK %s %s", p.name, instance.Id)
					wg.Done()
				}(instance)
				break
			}
		}
	}

	wg.Wait()
	return nil
}

func (p Provider) ShowRegions() error {
	regions, err := p.client.GetRegions(p.args)
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Printf("%s (%s)\n", region.Id, region.Country)
	}

	return nil
}

func New(providerName string, args data.ArgsBundle) (*Provider, error) {
	client, err := clients.New(providerName, args)
	if err != nil {
		return nil, err
	}

	return &Provider{
		name:   providerName,
		client: *client,
		args:   args,
	}, nil
}

func ListProviders() []string {
	list := make([]string, len(clients.AvailableProviders))
	i := 0
	for name := range clients.AvailableProviders {
		list[i] = name
		i++
	}
	return list
}

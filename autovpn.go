package main

import (
	"autovpn/helpers"
	"autovpn/openvpn"
	"autovpn/options"
	"autovpn/providers"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var usage = `
AutoVPN

Automatically provisions and de-provisions single-use VPN servers for one-shot VPN sessions.

Usage: autovpn <provider> <region>  Provision a VPN server at <provider> on <region> and connects to it
	   autovpn <provider> zombies   Lists all AutoVPN servers that should be destroyed at provider
       autovpn <provider> purge     Destroys all AutoVPN servers at provider
       autovpn <provider>           Lists all regions at <provider>
       
       autovpn providers            Lists all available providers
       autovpn zombies              Lists all AutoVPN servers that should be destroyed
       autovpn purge                Destroys all AutoVPN servers at all providers
       
       autovpn (-h | --help)        Shows further help and options
       autovpn --version            Shows version

Arguments:
  <provider>  VPS provider to use
  <region>    VPS provider region on which to create VPN endpoint

Options:
  -h --help    show this
  --version    show version`

var version = "DEVELOPMENT_BUILD"

func getProvider(name string) (providers.Provider, error) {
	provider, err := providers.New(name)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func showRegions(provider providers.Provider) error {
	regions, err := provider.GetRegions()
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Printf("%s (%s)\n", region.Id, region.Country)
	}

	return nil
}

func destroyServer(provider providers.Provider, server providers.Instance, key string) {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Destroying server", finish, exited)
	err := provider.DestroyServer(server, key)
	finish <- true
	<-exited
	if err != nil {
		panic(err)
	}
}

func purgeProvider(providerName string, config options.Config) error {
	provider, err := getProvider(providerName)
	if err != nil {
		return err
	}

	instances, err := provider.GetInstances(config)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, instance := range instances {
		isAutoVPNInstance := false
		for _, tag := range instance.Tags {
			if tag == providers.InstanceTag {
				isAutoVPNInstance = true
				break
			}
		}

		if !isAutoVPNInstance {
			continue
		}

		wg.Add(1)
		go func(instance providers.Instance) {
			log.Printf("Purging %s %s...", providerName, instance.Id)
			err := provider.DestroyServer(instance, config.Providers[providerName].Key)
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

func purgeAll(config options.Config) error {
	for name := range config.Providers {
		err := purgeProvider(name, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func listZombies(providerName string, config options.Config) error {
	provider, err := getProvider(providerName)
	if err != nil {
		return err
	}

	instances, err := provider.GetInstances(config)
	if err != nil {
		return err
	}

	fmt.Printf("--- %s: %d ---\n", providerName, len(instances))
	for _, instance := range instances {
		fmt.Printf("%s %s\n", providerName, instance.Id)
	}

	return nil
}

func listAllZombies(config options.Config) error {
	for name := range config.Providers {
		err := listZombies(name, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeOvpnConfig(ovpnConfig *string) {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint(fmt.Sprintf("Removing %s", *ovpnConfig), finish, exited)
	_ = os.Remove(*ovpnConfig)
	finish <- true
	<-exited
}

func provisionAndConnect(provider providers.Provider, arguments options.Arguments, config options.Config) error {
	if !arguments.NoAdminCheck {
		isAdmin, err := helpers.IsAdmin()
		if err != nil {
			return err
		}
		if !isAdmin {
			return fmt.Errorf("Root/Administrator privileges required")
		}
	}

	exe := openvpn.GetExecutable(config.Overrides.OpenvpnExe)
	if isInstalled := openvpn.IsInstalled(exe); !isInstalled {
		return fmt.Errorf("couldn't find OpenVPN exe (%s). OpenVPN must be installed", exe)
	}

	key := config.Providers[arguments.Provider].Key
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Creating instance", finish, exited)
	instance, err := provider.CreateServer(arguments, config)
	finish <- true
	<-exited
	if err != nil {
		if instance != nil {
			destroyServer(provider, *instance, key)
		}
		return err
	}
	defer destroyServer(provider, *instance, key)

	go helpers.WaitPrint("Starting instance", finish, exited)
	err = provider.AwaitProvisioning(*instance, key)
	finish <- true
	<-exited
	if err != nil {
		return err
	}

	go helpers.WaitPrint("Installing OpenVPN Server", finish, exited)
	ovpnConfig, err := openvpn.Install(*instance, config.Agent.ScriptUrl)
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

func main() {
	arguments := options.ParseArguments(os.Args)
	var configPath string

	if arguments.DebugMode {
		configPath = "./config.yml"
	} else {
		exe, err := os.Executable()
		if err != nil {
			log.Fatalln(err)
		}
		configPath = filepath.Dir(exe) + "/config.yml"
	}

	config, err := options.ReadConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if arguments.ShowHelp {
		fmt.Println(usage)
		os.Exit(0)

	} else if arguments.ShowVersion {
		fmt.Println(version)
		os.Exit(0)

	} else if arguments.ShowProviders {
		for _, provider := range providers.ListProviders() {
			fmt.Println(provider)
		}
		os.Exit(0)

	} else if arguments.Purge && len(arguments.Provider) == 0 {
		err = purgeAll(*config)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)

	} else if arguments.ListZombies && len(arguments.Provider) == 0 {
		err = listAllZombies(*config)
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	provider, err := getProvider(arguments.Provider)
	if err != nil {
		log.Fatalln(err)
	}

	if arguments.ShowRegions {
		err = showRegions(provider)
		if err != nil {
			log.Fatalln(err)
		}

	} else if arguments.Purge {
		err = purgeProvider(arguments.Provider, *config)
		if err != nil {
			log.Fatalln(err)
		}

	} else if arguments.ListZombies {
		err = listZombies(arguments.Provider, *config)
		if err != nil {
			log.Fatalln(err)
		}

	} else {
		err = provisionAndConnect(provider, arguments, *config)
		if err != nil {
			log.Fatalln(err)
		}
	}

	os.Exit(0)
}

package main

import (
	"autovpn/helpers"
	"autovpn/options"
	"autovpn/providers"
	"fmt"
	"os"
)

var usage = `
AutoVPN

Automatically provisions and de-provisions single-use VPN servers for one-shot VPN sessions.

Usage: autovpn <provider> <region>
       autovpn <provider> (--silent | -s)
       autovpn providers
       autovpn (-h | --help)
       autovpn --version

Commands:
  <provider> <region>  create and connect to VPN endpoint at <provider> on <region>
  <provider>           list available regions for <provider>
  providers            list available providers

Arguments:
  <provider>  VPS provider to use
  <region>    VPS provider region on which to create VPN endpoint

Options:
  -s --silent  hush logging to stdout
  -h --help    show this
  --version    show version`

var version = "DEVELOPMENT_BUILD"

func showRegions(provider providers.Provider, silently bool) error {
	finish := make(chan bool)
	exited := make(chan bool)

	if !silently {
		go helpers.WaitPrint("Downloading regions", finish, exited)
	}
	regions, err := provider.GetRegions()
	if !silently {
		finish <- true
		<-exited
	}
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Printf("%s (%s)\n", region.Id, region.Country)
	}

	return nil
}

func destroyServer(provider providers.Provider, server providers.Instance, key string) error {
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Destroying server", finish, exited)
	err := provider.DestroyServer(server, key)
	finish <- true
	<-exited
	if err != nil {
		return err
	}

	return nil
}

func provisionAndConnect(provider providers.Provider, arguments options.Arguments, config options.Config) error {
	var server *providers.Instance = nil
	key := config.Providers[arguments.Provider].Key
	finish := make(chan bool)
	exited := make(chan bool)

	go helpers.WaitPrint("Creating server", finish, exited)
	server, err := provider.CreateServer(arguments, config)
	finish <- true
	<-exited
	if err != nil {
		if server != nil {
			err := destroyServer(provider, *server, key)
			if err != nil {
				return err
			}
		}
		return err
	}

	go helpers.WaitPrint("Provisioning server", finish, exited)
	err = provider.AwaitProvisioning(*server, key)
	finish <- true
	<-exited
	if err != nil {
		err := destroyServer(provider, *server, key)
		if err != nil {
			return err
		}
		return err
	}

	err = destroyServer(provider, *server, key)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	arguments, err := options.ParseArguments(os.Args)
	if err != nil {
		panic(err)
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
	}

	provider, err := providers.New(arguments.Provider)
	if err != nil {
		panic(err)
	}

	if arguments.ShowRegions {
		err := showRegions(provider, arguments.Silent)
		if err != nil {
			panic(err)
		}
	} else {
		conf, err := options.ReadConfig("./config.yml")
		if err != nil {
			panic(err)
		}

		err = provisionAndConnect(provider, arguments, *conf)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(0)
}

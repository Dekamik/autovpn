package main

import (
	"autovpn/config"
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
	regions, err := provider.GetRegions(silently)
	if err != nil {
		return err
	}

	for _, region := range regions {
		fmt.Printf("%s (%s)\n", region.Id, region.Country)
	}

	return nil
}

func provisionAndConnect(provider providers.Provider, arguments config.Arguments, yamlConfig config.YamlConfig) error {
	providerConfig, err := providers.GetProviderConfig(yamlConfig, arguments.Provider)
	if err != nil {
		return err
	}

	server, err := provider.CreateServer(arguments, providerConfig)
	if err != nil {
		return err
	}

	err = provider.DestroyServer(*server, "")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	arguments, err := config.ParseArguments(os.Args)
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
		for _, provider := range providers.AvailableProviders {
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
		conf, err := config.ReadYamlConfig("./config.yml")
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

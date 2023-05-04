package main

import (
	"autovpn/data"
	"autovpn/providers"
	"fmt"
	"os"
	"path/filepath"
)

var usage = `
AutoVPN

Automatically provisions and de-provisions single-use VPN servers for one-shot VPN sessions.

Usage: autovpn <provider> <region>  Provision a VPN server at <provider> on <region> and connects to it
       autovpn <provider>           Lists all regions at <provider> (e.g: linode)
       autovpn providers            Lists all available providers

       autovpn [provider] zombies  Lists all AutoVPN servers that should be destroyed. Lists all zombies across all providers if provider is omitted
       autovpn [provider] purge    Destroys all AutoVPN servers at all providers. Destroys all AutoVPN servers at all providers if provider is omitted
       
       autovpn (-h | --help)        Shows further help and options
       autovpn --version            Shows version

Arguments:
  <provider>  VPS provider to use
  <region>    VPS provider region on which to create VPN endpoint

Options:
  -h --help    show this
  --version    show version`

var version = "DEVELOPMENT_BUILD"

// purgeAll destroys all AutoVPN servers on every provider.
func purgeAll(args data.ArgsBundle) error {
	for _, providerName := range providers.ListProviders() {
		provider, err := providers.New(providerName, args)
		if err != nil {
			return err
		}

		err = provider.Purge()
		if err != nil {
			return err
		}
	}
	return nil
}

// listAllZombies goes through every provider and checks for zombies.
func listAllZombies(args data.ArgsBundle) error {
	for _, providerName := range providers.ListProviders() {
		provider, err := providers.New(providerName, args)
		if err != nil {
			return err
		}

		err = provider.ListZombies()
		if err != nil {
			return err
		}
	}
	return nil
}

func getArgs() (*data.ArgsBundle, *providers.Provider, error) {
	arguments, err := data.ParseArguments(os.Args)
	if err != nil {
		return nil, nil, err
	}

	var configPath string
	if arguments.DebugMode {
		configPath = "./config.yml"
	} else {
		exe, err := os.Executable()
		if err != nil {
			return nil, nil, err
		}
		configPath = filepath.Dir(exe) + "/config.yml"
	}

	config, err := data.ReadConfig(configPath)
	if err != nil {
		return nil, nil, err
	}

	args := &data.ArgsBundle{
		Config:    *config,
		Arguments: *arguments,
	}

	var provider *providers.Provider
	if len(arguments.Provider) != 0 {
		provider, err = providers.New(arguments.Provider, *args)
		if err != nil {
			return nil, nil, err
		}
	}

	return args, provider, nil
}

func main() {
	args, provider, err := getArgs()
	if err != nil {
		fmt.Printf("\n%s", err)
		os.Exit(1)
	}

	exitCode := 0

	switch args.Arguments.Command {

	case data.ListProviders:
		for _, provider := range providers.ListProviders() {
			fmt.Println(provider)
		}

	case data.ListRegions:
		err = provider.ShowRegions()
		if err != nil {
			fmt.Printf("\n%s", err)
			exitCode = 1
		}

	case data.ListZombies:
		if provider == nil {
			err = listAllZombies(*args)
			if err != nil {
				fmt.Printf("\n%s", err)
				exitCode = 1
			}
		} else {
			err = provider.ListZombies()
			if err != nil {
				fmt.Printf("\n%s", err)
				exitCode = 1
			}
		}

	case data.Purge:
		if provider == nil {
			err = purgeAll(*args)
			if err != nil {
				fmt.Printf("\n%s", err)
				exitCode = 1
			}
		} else {
			err = provider.Purge()
			if err != nil {
				fmt.Printf("\n%s", err)
				exitCode = 1
			}
		}

	case data.Version:
		fmt.Println(version)

	case data.Usage:
		fmt.Println(usage)

	default:
		err = provider.Connect()
		if err != nil {
			fmt.Printf("\n%s", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

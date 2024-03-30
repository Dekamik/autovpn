package main

import (
	"autovpn/internal/data"
	"autovpn/internal/providers"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
)

var usage = `
AutoVPN

Automatically provisions and de-provisions single-use VPN servers for one-shot
VPN sessions.

Example: autovpn linode se-sto      Provisions a VPN server at Linode's
                                    datacenter in Stockholm, Sweden

Usage:	autovpn <provider> <region> Provision a VPN server at the specified
                                    provider on the specified region and then
                                    connects to it

    	autovpn <provider> list     Lists all region slugs at the provider

    	autovpn list                Lists all available providers

    	autovpn <provider> zombies  Lists all AutoVPN servers that should be
                                    destroyed at the provider

    	autovpn <provider> purge    Destroys all AutoVPN servers at the provider

    	autovpn (--help)            Shows further help and options

    	autovpn --version           Shows version

Arguments:
	<provider>  VPS provider to use
	<region>    VPS provider region on which to create the VPN endpoint

Options:
	-c <path>   specify a configuration file

	--help      show this
	--version   show version`

var version = "DEVELOPMENT_BUILD"

var ErrConfigNotFound = errors.New("config not found")

func setup() (*data.ArgsBundle, *providers.Provider, error) {
	arguments, err := data.ParseArguments()
	if err != nil {
		return nil, nil, err
	}

	var configPath string = ""
	if arguments.PredefinedConfigPath != "" {
		if _, err = os.Stat(arguments.PredefinedConfigPath); err != nil {
			log.Printf("PredefinedConfigPath = %s", arguments.PredefinedConfigPath)
			return nil, nil, ErrConfigNotFound
		}

		configPath = arguments.PredefinedConfigPath

	} else {
		var home string

		if username := os.Getenv("SUDO_USER"); username != "" {
			u, err := user.Lookup(username)
			if err != nil {
				return nil, nil, err
			}
			home = u.HomeDir
		} else {
			u, err := user.Current()
			if err != nil {
				return nil, nil, err
			}
			home = u.HomeDir
		}

		configPaths := []string {
			"./.autovpn.yml",
			home + "/.autovpn.yml",
		}

		for _, path := range configPaths {
			if _, err = os.Stat(path); err == nil {
				configPath = path
				break
			}
		}

		if configPath == "" {
			return nil, nil, ErrConfigNotFound
		}
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
	args, provider, err := setup()
	if err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	var exitCode int = 0

	switch args.Arguments.Command {

	case data.ListArgs:
		if args.Arguments.Provider == "" {
			for _, p := range providers.ListProviders() {
				fmt.Println(p)
			}
		} else {
			err = provider.ShowRegions()
			if err != nil {
				fmt.Printf("\n%s\n", err)
				exitCode = 1
			}
		}

	case data.ListZombies:
		err = provider.ListZombies()
		if err != nil {
			fmt.Printf("\n%s\n", err)
			exitCode = 1
		}

	case data.Purge:
		err = provider.Purge()
		if err != nil {
			fmt.Printf("\n%s\n", err)
			exitCode = 1
		}

	case data.Version:
		fmt.Println(version)

	case data.Usage:
		fmt.Println(usage)

	default:
		err = provider.Connect()
		if err != nil {
			fmt.Printf("\n%s\n", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

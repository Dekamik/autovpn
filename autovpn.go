package main

import (
	"autovpn/options"
	"autovpn/providers"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

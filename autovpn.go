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
)

var usage = `
AutoVPN

Automatically provisions and de-provisions single-use VPN servers for one-shot VPN sessions.

Usage: autovpn <provider> <region>
       autovpn <provider>
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
  -h --help    show this
  --version    show version`

var version = "DEVELOPMENT_BUILD"

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
		log.Fatalln(err)
	}

	if arguments.ShowRegions {
		err := showRegions(provider)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
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

		conf, err := options.ReadConfig(configPath)
		if err != nil {
			log.Fatalln(err)
		}

		err = provisionAndConnect(provider, arguments, *conf)
		if err != nil {
			log.Fatalln(err)
		}
	}

	os.Exit(0)
}

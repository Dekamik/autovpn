package data

import (
	"errors"
	"flag"
)

var ErrInvalidFirstArgumentRegions = errors.New(`invalid first argument "regions"`)
var ErrInvalidSecondArgumentProviders = errors.New(`invalid second argument "providers"`)

type Command int64

const (
	Default Command = iota
	ListArgs
	ListZombies
	Purge

	Version
	Usage
)

type Arguments struct {
	Command    Command
	Provider   string
	Region     string

	PredefinedConfigPath string

	DebugMode    bool
	NoAdminCheck bool
}

func getCommand(arg string) Command {
	switch arg {
	case "list":
		return ListArgs
	case "purge":
		return Purge
	case "zombies":
		return ListZombies
	default:
		return Default
	}
}

func ParseArguments() (*Arguments, error) {
	var configPath = flag.String("c", "", "define which config to use")
	var debugMode = flag.Bool("debug", false, "run application in debug mode")
	var noAdminCheck = flag.Bool("no-admin-check", false, "run even when not running as a privileged user")
	var help = flag.Bool("help", false, "display usage")
	var version = flag.Bool("version", false, "display version")

	flag.Parse()

	var arg1 = flag.Arg(0)
	var arg2 = flag.Arg(1)

	if *help || arg1 == "" {
		return &Arguments{Command: Usage}, nil
	}
	if *version {
		return &Arguments{Command: Version}, nil
	}

	var command Command = Default
	var provider string
	var region string
	var err error

	if arg2 == "" {
		command = getCommand(arg1)
		if command == Default {
			command = Usage
		}
	} else {
		command = getCommand(arg2)
		provider = arg1
		if command == Default {
			region = arg2
		}
	}

	if err != nil {
		return nil, err
	}

	arguments := Arguments{
		Command:              command,
		Provider:             provider,
		Region:               region,
		PredefinedConfigPath: *configPath,
		DebugMode:            *debugMode,
		NoAdminCheck:         *noAdminCheck,
	}

	return &arguments, nil
}

package data

import (
	"flag"
)

type Command int64

const (
	Default Command = iota
	ListProviders
	ListRegions
	ListZombies
	Purge

	Version
	Usage
)

type Arguments struct {
	Command  Command
	Provider string
	Region   string

	DebugMode    bool
	NoAdminCheck bool
}

func ParseArguments(argv []string) (*Arguments, error) {
	if *flag.Bool("help", false, "display usage") {
		return &Arguments{Command: Usage}, nil
	}
	if *flag.Bool("version", false, "display version") {
		return &Arguments{Command: Version}, nil
	}

	arguments := Arguments{}
	arguments.DebugMode = *flag.Bool("debug", false, "run application in debug mode")
	arguments.NoAdminCheck = *flag.Bool("no-admin-check", false, "run even when not running as a privileged user")

	var arg1 = flag.Arg(1)
	switch arg1 {
	case "providers":
		arguments.Command = ListProviders
	case "purge":
		arguments.Command = Purge
	case "zombies":
		arguments.Command = ListZombies
	default:
		if len(arg1) != 0 {
			arguments.Provider = arg1
		} else {
			arguments.Command = ListRegions
		}
	}

	arguments.Region = flag.Arg(2)

	return &arguments, nil
}

package data

import "fmt"

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
	arguments := Arguments{}

	if len(argv) == 1 {
		return &Arguments{Command: Usage}, nil
	}

	for _, arg := range argv[1:] {
		switch arg {

		case "--help":
		case "-h":
			return &Arguments{Command: Usage}, nil

		case "--version":
			return &Arguments{Command: Version}, nil

		case "--debug":
			arguments.DebugMode = true

		case "--no-admin-check":
			arguments.NoAdminCheck = true

		case "providers":
			arguments.Command = ListProviders

		case "purge":
			arguments.Command = Purge

		case "zombies":
			arguments.Command = ListZombies

		default:
			if len(arguments.Provider) == 0 {
				arguments.Provider = arg

			} else if len(arguments.Region) == 0 {
				arguments.Region = arg

			} else {
				return nil, fmt.Errorf("unexpected third argument: %s", arg)
			}
		}
	}

	if arguments.Command == Default && len(arguments.Region) == 0 {
		arguments.Command = ListRegions
	}

	return &arguments, nil
}

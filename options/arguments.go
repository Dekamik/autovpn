package options

type Arguments struct {
	Provider    string
	Region      string
	Purge       bool
	ListZombies bool

	DebugMode    bool
	NoAdminCheck bool

	ShowProviders bool
	ShowRegions   bool

	ShowHelp    bool
	ShowVersion bool
}

func ParseArguments(argv []string) Arguments {
	arguments := Arguments{}

	if len(argv) == 1 {
		return Arguments{ShowHelp: true}
	}

	for _, arg := range argv[1:] {
		switch arg {

		case "--help":
		case "-h":
			return Arguments{ShowHelp: true}

		case "--version":
			return Arguments{ShowVersion: true}

		case "--debug":
			arguments.DebugMode = true

		case "--no-admin-check":
			arguments.NoAdminCheck = true

		case "providers":
			return Arguments{ShowProviders: true}

		case "purge":
			arguments.Purge = true

		case "zombies":
			arguments.ListZombies = true

		default:
			if len(arguments.Provider) == 0 {
				arguments.Provider = arg
			} else {
				arguments.Region = arg
			}
		}
	}

	return arguments
}

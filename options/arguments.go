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
			arguments.ShowHelp = true
			break

		case "--version":
			arguments.ShowVersion = true
			break

		case "--debug":
			arguments.DebugMode = true
			break

		case "--no-admin-check":
			arguments.NoAdminCheck = true
			break

		case "providers":
			arguments.ShowProviders = true
			break

		case "purge":
			arguments.Purge = true
			break

		case "zombies":
			arguments.ListZombies = true
			break

		default:
			if len(arguments.Provider) == 0 {
				arguments.Provider = arg
			} else {
				arguments.Region = arg
			}
			break
		}
	}

	if len(arguments.Region) == 0 {
		arguments.ShowRegions = true
	}

	return arguments
}

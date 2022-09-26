package options

type Arguments struct {
	Provider string
	Region   string

	ShowProviders bool
	ShowRegions   bool
	Silent        bool

	ShowHelp    bool
	ShowVersion bool
}

func ParseArguments(argv []string) (Arguments, error) {
	var provider string

	if len(argv) == 1 {
		return Arguments{ShowHelp: true}, nil
	}

	switch argv[1] {

	case "--help":
	case "-h":
		return Arguments{ShowHelp: true}, nil

	case "--version":
		return Arguments{ShowVersion: true}, nil

	case "providers":
		return Arguments{ShowProviders: true}, nil

	default:
		provider = argv[1]
	}

	if len(argv) <= 2 {
		return Arguments{Provider: provider, ShowRegions: true}, nil
	} else if argv[2] == "-s" || argv[2] == "--silent" {
		return Arguments{Provider: provider, ShowRegions: true, Silent: true}, nil
	}
	return Arguments{Provider: provider, Region: argv[2]}, nil
}

package data

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent struct {
		ScriptUrl string `yaml:"script_url"`
	}
	Overrides struct {
		OpenvpnExe string `yaml:"openvpn_exe"`
		RootPass   string `yaml:"root_pass"`
	}
	Providers map[string]struct {
		Image    string
		Key      string
		TypeSlug string `yaml:"type_slug"`
	}
}

func ReadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		log.Printf("error when reading config in file %s", path)
		return nil, err
	}

	return conf, nil
}

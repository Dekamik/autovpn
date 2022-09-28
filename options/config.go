package options

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent struct {
		ScriptUrl string `yaml:"script_url"`
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
		return nil, fmt.Errorf("in file %q: %w", path, err)
	}

	return conf, nil
}

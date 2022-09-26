package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Agent struct {
		ScriptUrl string `yaml:"script_url"`
	}
	Openvpn struct {
		Path struct {
			Darwin  string
			Linux   string
			Windows string
		}
	}
	Providers map[string]interface{}
}

func ReadYamlConfig(path string) (*YamlConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &YamlConfig{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", path, err)
	}

	return conf, nil
}

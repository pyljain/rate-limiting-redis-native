package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Defaults ConfigDefaults `yaml:"defaults"`
	Rules    []ConfigRule   `yaml:"rules"`
}

type ConfigDefaults struct {
	CSI    int64         `yaml:"csi"`
	Models []ConfigModel `yaml:"models"`
}

type ConfigModel struct {
	Name  string `yaml:"name"`
	Value int64  `yaml:"value"`
}

type ConfigRule struct {
	CSI    string        `yaml:"csi"`
	Limit  int64         `yaml:"limit"`
	Models []ConfigModel `yaml:"models"`
}

func Load() (*Config, error) {
	contents, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

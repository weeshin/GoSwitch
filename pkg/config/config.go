package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Channels []ChannelConfig `yaml:"channels"`
}

type ServerConfig struct {
	Port        int    `yaml:"port"`
	IP          string `yaml:"ip"`
	ReadTimeout int    `yaml:"read_timeout"`
}

type ChannelConfig struct {
	Name              string `yaml:"name"`
	IP                string `yaml:"ip"`
	Port              int    `yaml:"port"`
	Enabled           bool   `yaml:"enabled"`
	ReconnectInterval int    `yaml:"reconnect_interval"`
}

func LoadAppConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

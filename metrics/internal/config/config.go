package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Env string `yaml:"env"`

	Port int `yaml:"port"`

	Broker Broker `yaml:"broker"`
}

type Broker struct {
	Addr   string `yaml:"addr"`
	Topics Topics `yaml:"topics"`
}

type Topics struct {
	CarsharingResponseTime string `yaml:"responseTime"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("path to config file is not provided")
	}

	if _, err := os.Stat(path); err != nil {
		panic("can not find config file: " + path)
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}

	var cfg Config
	if err = yaml.Unmarshal(configFile, &cfg); err != nil {
		panic("failed to unmarshal config file")
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "c", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}

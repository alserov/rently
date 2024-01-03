package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port int
	Env  string

	DB Postgres

	Broker Broker
}

type Postgres struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
}

type Broker struct {
	Addr    string `yaml:"addr"`
	Metrics struct {
		Topics struct {
			DecreaseActiveRentsAmount string `yaml:"decreaseActiveRentsAmount"`
			IncreaseActiveRentsAmount string `yaml:"increaseActiveRentsAmount"`
			IncreaseRentsCancel       string `yaml:"increaseRentsCancel"`
			NotifyBrandDemand         string `yaml:"notifyBrandDemand"`
		} `yaml:"topics"`
	} `yaml:"metrics"`
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

package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Port int

	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`

	Env string

	Services Services
}

type Services struct {
	CarAddr  string `yaml:"carAddr"`
	UserAddr string `yaml:"userAddr"`
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
		panic("failed to unmarshal config file: " + err.Error())
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

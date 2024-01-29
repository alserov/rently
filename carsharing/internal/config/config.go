package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port     int
	Env      string
	DB       Postgres
	Cache    Cache
	Services Services
	Broker   Broker
}

type Cache struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type Services struct {
	Payment struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"payment"`
}

type Postgres struct {
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
}

func (p *Postgres) GetDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.Name)
}

type Broker struct {
	Addr   string `yaml:"addr"`
	Topics struct {
		Metrics      Metrics `yaml:"metrics"`
		Notification struct {
			RentCreated  string `yaml:"rentCreated"`
			RentCanceled string `yaml:"rentCanceled"`
		}
	} `yaml:"topics"`
}

type Metrics struct {
	DecreaseActiveRentsAmount string `yaml:"decreaseActiveRentsAmount"`
	IncreaseActiveRentsAmount string `yaml:"increaseActiveRentsAmount"`
	IncreaseRentsCancel       string `yaml:"increaseRentsCancel"`
	NotifyBrandDemand         string `yaml:"notifyBrandDemand"`
	ResponseTime              string `yaml:"responseTime"`
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

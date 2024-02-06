package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Env string `yaml:"env"`

	Port int `yaml:"port"`

	DB Mysql `yaml:"db"`

	Broker Broker `yaml:"broker"`
}

type Mysql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

func (mysql *Mysql) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Name)
}

type Broker struct {
	Rabbit struct {
		Addr   string `yaml:"addr"`
		Topics Topics `yaml:"topics"`
	} `yaml:"rabbit"`
}

type Topics struct {
	Email string `yaml:"email"`
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

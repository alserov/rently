package config

type Config interface {}

func NewConfig() Config {
	return &config{}
}

type config struct {}
			
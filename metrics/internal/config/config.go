package config

type Config struct {
	Port int `yaml:"port"`
}

type Broker struct {
	Addr   string `yaml:"addr"`
	Topics struct {
		RentAmount string
	} `yaml:"topics"`
}

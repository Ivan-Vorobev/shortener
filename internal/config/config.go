package config

import (
	"flag"
	"os"
)

type Configuration struct {
	ServerAddress string
	BaseURL       string
}

func NewDefaultConfig() *Configuration {
	return &Configuration{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}
}

func LoadConfiguration() *Configuration {
	config := NewDefaultConfig()

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "address to run HTTP server")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base URL for short links")
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.ServerAddress = addr
	}

	if bastURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseURL = bastURL
	}

	return config
}

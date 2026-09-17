package config

import (
	"flag"
	"os"
)

type Configuration struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
}

func NewDefaultConfig() *Configuration {
	return &Configuration{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "link-storage.json",
	}
}

func LoadConfiguration() *Configuration {
	config := NewDefaultConfig()

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "address to run HTTP server")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base URL for short links")
	flag.StringVar(&config.FileStoragePath, "f", config.FileStoragePath, "file storage path")
	flag.Parse()

	if addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		config.ServerAddress = addr
	}

	if bastURL, ok := os.LookupEnv("BASE_URL"); ok {
		config.BaseURL = bastURL
	}

	if fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		config.FileStoragePath = fileStoragePath
	}

	return config
}

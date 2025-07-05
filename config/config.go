package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	Name         string `yaml:"name"`
	UpstreamPath string `yaml:"upstream_path"`
}

type Cors struct {
	AllowOrigins string `yaml:"allow_origins"`
	AllowMethods string `yaml:"allow_methods"`
	AllowHeaders string `yaml:"allow_headers"`
}

type Config struct {
	Port     string             `yaml:"port"`
	Services map[string]Service `yaml:"services"`
	Cors     Cors               `yaml:"cors"`
}

func GetConfig(configurationFilePath string) *Config {
	var c Config

	f, err := os.ReadFile(configurationFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	return &c
}

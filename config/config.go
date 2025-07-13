package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Service struct {
	UpstreamPath string `yaml:"upstream_path"`
}

type Cors struct {
	AllowOrigins string `yaml:"allow_origins"`
	AllowMethods string `yaml:"allow_methods"`
	AllowHeaders string `yaml:"allow_headers"`
}

type PerIpAddress struct {
	Enabled                bool `yaml:"enabled"`
	CleanupIntervalSeconds int  `yaml:"cleanup_interval_seconds"`
}

type RateLimiting struct {
	MaxCapacity  int64        `yaml:"max_capacity"`
	PerIpAddress PerIpAddress `yaml:"per_ip_address"`
}

type Auth struct {
	UseJwt bool `yaml:"use_jwt"`
}

type Config struct {
	Port         string             `yaml:"port"`
	Services     map[string]Service `yaml:"services"`
	Cors         Cors               `yaml:"cors"`
	RateLimiting RateLimiting       `yaml:"rate_limiting"`
	Auth         Auth               `yaml:"auth"`
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

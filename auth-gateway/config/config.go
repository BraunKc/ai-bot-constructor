package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Env string

type LoggerOutputType string

type LoggerLevel string

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AuthServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConfig struct {
	AuthService AuthServiceConfig `yaml:"auth-service"`
}

type LoggerConfig struct {
	Service    string           `yaml:"service"`
	OutputType LoggerOutputType `yaml:"output-type"`
	Level      LoggerLevel      `yaml:"level"`
}

type Config struct {
	Env    Env          `yaml:"env"`
	GRPC   GRPCConfig   `yaml:"grpc"`
	Logger LoggerConfig `yaml:"logger"`
}

// TODO: write validation
func New(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}
	file = []byte(os.ExpandEnv(string(file)))

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

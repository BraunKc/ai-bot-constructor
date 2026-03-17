package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type Env string

type LoggerOutputType string

type LoggerLevel string

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type OrchestratorServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConfig struct {
	OrchestratorService OrchestratorServiceConfig `yaml:"orchestrator-service"`
}

type LoggerConfig struct {
	Service    string           `yaml:"service"`
	OutputType LoggerOutputType `yaml:"output-type"`
	Level      LoggerLevel      `yaml:"level"`
}

type Config struct {
	Env    Env          `yaml:"env"`
	HTTP   HTTPConfig   `yaml:"http"`
	GRPC   GRPCConfig   `yaml:"grpc"`
	Logger LoggerConfig `yaml:"logger"`
}

func New(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	file, err := os.ReadFile(os.Getenv("CONFIG_PATH"))
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

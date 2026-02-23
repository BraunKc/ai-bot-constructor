package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Env string

type LoggerOutputType string

type LoggerLevel string

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type LoggerConfig struct {
	Service    string           `yaml:"service"`
	OutputType LoggerOutputType `yaml:"output-type"`
	Level      LoggerLevel      `yaml:"level"`
}

type Config struct {
	GRPC   GRPCConfig   `yaml:"grpc"`
	DB     DBConfig     `yaml:"db"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Logger LoggerConfig `yaml:"logger"`
}

// TODO: write validate
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

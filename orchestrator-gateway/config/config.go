package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var (
	ErrEmptyConfigPath = errors.New("empty CONFIG_PATH")

	ErrInvalidEnv = errors.New("invalid env")

	ErrEmptyHTTPHost = errors.New("empty HTTP_HOST")
	ErrEmptyHTTPPort = errors.New("empty HTTP_PORT")

	ErrEmptyGRPCHost = errors.New("empty GRPC_ORCHESTRATOR_HOST")
	ErrEmptyGRPCPort = errors.New("empty GRPC_ORCHESTRATOR_PORT")

	ErrEmptyLoggerService      = errors.New("invalid logger.service")
	ErrInvalidLoggerOutputType = errors.New("invalid logger.output-type")
	ErrInvalidLoggerLevel      = errors.New("invalid logger.level")
)

const (
	Develop    Env = "develop"
	Production Env = "production"

	Console LoggerOutputType = "console"
	File    LoggerOutputType = "file"
	Both    LoggerOutputType = "both"

	Debug LoggerLevel = "debug"
	Info  LoggerLevel = "info"
	Warn  LoggerLevel = "warn"
	Error LoggerLevel = "error"
)

type Env string

func (e Env) Validate() error {
	switch e {
	case Develop, Production:
		return nil
	default:
		return ErrInvalidEnv
	}
}

type LoggerOutputType string

func (t LoggerOutputType) Validate() error {
	switch t {
	case Console, File, Both:
		return nil
	default:
		return ErrInvalidLoggerOutputType
	}
}

type LoggerLevel string

func (l LoggerLevel) Validate() error {
	switch l {
	case Debug, Info, Warn, Error:
		return nil
	default:
		return ErrInvalidLoggerLevel
	}
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (h *HTTPConfig) Validate() error {
	switch {
	case h.Host == "":
		return ErrEmptyHTTPHost
	case h.Port == "":
		return ErrEmptyHTTPPort
	}
	return nil
}

type OrchestratorServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (o *OrchestratorServiceConfig) Validate() error {
	switch {
	case o.Host == "":
		return ErrEmptyGRPCHost
	case o.Port == "":
		return ErrEmptyGRPCPort
	}
	return nil
}

type GRPCConfig struct {
	OrchestratorService OrchestratorServiceConfig `yaml:"orchestrator-service"`
}

func (g *GRPCConfig) Validate() error {
	return g.OrchestratorService.Validate()
}

type LoggerConfig struct {
	Service    string           `yaml:"service"`
	OutputType LoggerOutputType `yaml:"output-type"`
	Level      LoggerLevel      `yaml:"level"`
}

func (l *LoggerConfig) Validate() error {
	if l.Service == "" {
		return ErrEmptyLoggerService
	}

	if err := l.OutputType.Validate(); err != nil {
		return err
	}

	if err := l.Level.Validate(); err != nil {
		return err
	}

	return nil
}

type Config struct {
	Env    Env          `yaml:"env"`
	HTTP   HTTPConfig   `yaml:"http"`
	GRPC   GRPCConfig   `yaml:"grpc"`
	Logger LoggerConfig `yaml:"logger"`
}

func New(path string) (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, ErrEmptyConfigPath
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}
	file = []byte(os.ExpandEnv(string(file)))

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if err := c.Env.Validate(); err != nil {
		return err
	}

	if err := c.HTTP.Validate(); err != nil {
		return err
	}

	if err := c.GRPC.Validate(); err != nil {
		return err
	}

	if err := c.Logger.Validate(); err != nil {
		return err
	}

	return nil
}

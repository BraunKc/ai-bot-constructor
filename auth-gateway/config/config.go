package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
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

var (
	ErrInvalidEnv = errors.New("invalid env")

	ErrEmptyHTTPHost = errors.New("empty HTTP_HOST")
	ErrEmptyHTTPPort = errors.New("empty HTTP_PORT")

	ErrEmptyAuthServiceHost = errors.New("empty GRPC_AUTH_HOST")
	ErrEmptyAuthServicePort = errors.New("empty GRPC_AUTH_PORT")

	ErrEmptyLoggerService      = errors.New("empty logger.service")
	ErrInvalidLoggerOutputType = errors.New("invalid logger.output-type")
	ErrInvalidLoggerLevel      = errors.New("invalid logger.level")
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

func (c *HTTPConfig) Validate() error {
	if c.Host == "" {
		return ErrEmptyHTTPHost
	}
	if c.Port == "" {
		return ErrEmptyHTTPPort
	}

	return nil
}

type AuthServiceConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (c *AuthServiceConfig) Validate() error {
	if c.Host == "" {
		return ErrEmptyAuthServiceHost
	}
	if c.Port == "" {
		return ErrEmptyAuthServicePort
	}

	return nil
}

type GRPCConfig struct {
	AuthService AuthServiceConfig `yaml:"auth-service"`
}

func (c *GRPCConfig) Validate() error {
	return c.AuthService.Validate()
}

type LoggerConfig struct {
	Service    string           `yaml:"service"`
	OutputType LoggerOutputType `yaml:"output-type"`
	Level      LoggerLevel      `yaml:"level"`
}

func (c *LoggerConfig) Validate() error {
	if c.Service == "" {
		return ErrEmptyLoggerService
	}
	if err := c.OutputType.Validate(); err != nil {
		return err
	}
	if err := c.Level.Validate(); err != nil {
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

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

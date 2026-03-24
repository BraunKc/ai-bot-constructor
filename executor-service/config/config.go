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

	ErrEmptyOpenRouterToken = errors.New("empty OPEN_ROUTER_TOKEN")

	ErrEmptyKafkaHost = errors.New("empty KAFKA_HOST")
	ErrEmptyKafkaPort = errors.New("empty KAFKA_PORT")

	ErrEmptyDBHost     = errors.New("empty DB_HOST")
	ErrEmptyDBPort     = errors.New("empty DB_PORT")
	ErrEmptyDBName     = errors.New("empty DB_NAME")
	ErrEmptyDBUser     = errors.New("empty DB_USER")
	ErrEmptyDBPassword = errors.New("empty DB_PASSWORD")

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

type KafkaConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (c *KafkaConfig) Validate() error {
	if c.Host == "" {
		return ErrEmptyKafkaHost
	}
	if c.Port == "" {
		return ErrEmptyKafkaPort
	}
	return nil
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (c *DBConfig) Validate() error {
	if c.Host == "" {
		return ErrEmptyDBHost
	}
	if c.Port == "" {
		return ErrEmptyDBPort
	}
	if c.Name == "" {
		return ErrEmptyDBName
	}
	if c.User == "" {
		return ErrEmptyDBUser
	}
	if c.Password == "" {
		return ErrEmptyDBPassword
	}
	return nil
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
	Env             Env          `yaml:"env"`
	OpenRouterToken string       `yaml:"open-router-token"`
	Kafka           KafkaConfig  `yaml:"kafka"`
	DB              DBConfig     `yaml:"db"`
	Logger          LoggerConfig `yaml:"logger"`
}

func (c *Config) Validate() error {
	if err := c.Env.Validate(); err != nil {
		return err
	}
	if c.OpenRouterToken == "" {
		return ErrEmptyOpenRouterToken
	}
	if err := c.Kafka.Validate(); err != nil {
		return err
	}
	if err := c.DB.Validate(); err != nil {
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

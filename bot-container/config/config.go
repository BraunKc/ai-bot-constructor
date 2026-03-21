package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

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

	ErrEmptyTelegramToken = errors.New("empty TELEGRAM_BOT_TOKEN")

	ErrEmptyLoggerService      = errors.New("empty logger.service")
	ErrInvalidLoggerOutputType = errors.New("invalid logger.output-type")
	ErrInvalidLoggerLevel      = errors.New("invalid logger.level")

	ErrEmptySystemPrompt = errors.New("empty SYSTEM_PROMPT")
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

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type OpenRouterConfig struct {
	Token string `yaml:"token"`
	Model string `yaml:"model"`
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
	Env          Env              `yaml:"env"`
	Telegram     TelegramConfig   `yaml:"tg"`
	OpenRouter   OpenRouterConfig `yaml:"open-router"`
	Logger       LoggerConfig     `yaml:"logger"`
	SystemPrompt string           `yaml:"system-prompt"`
}

func (c *Config) Validate() error {
	if err := c.Env.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(c.Telegram.Token) == "" {
		return ErrEmptyTelegramToken
	}
	if err := c.Logger.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(c.SystemPrompt) == "" {
		return ErrEmptySystemPrompt
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

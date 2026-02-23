package botdomain

import (
	"errors"
	"strings"
)

var (
	ErrInvalidBotStatus     = errors.New("invalid bot_status")
	ErrBotNameMustBeLonger  = errors.New("bot_name must be longer")
	ErrBotNameMustBeShorter = errors.New("bot_name must be shorter")
	ErrInvalidApiKey        = errors.New("invalid api_key")
)

const (
	BotStatusUnknown Status = iota
	BotStatusCreating
	BotStatusStarting
	BotStatusRunning
	BotStatusStopping
	BotStatusStopped
	BotStatusRestarting
	BotStatusDeleting
	BotStatusError
)

type Status int32

func NewStatus(status int32) (Status, error) {
	s := Status(status)
	if s < BotStatusUnknown || s > BotStatusError {
		return 0, ErrInvalidBotStatus
	}

	return s, nil
}

func (s Status) String() string {
	switch s {
	case BotStatusCreating:
		return "creating"
	case BotStatusStarting:
		return "starting"
	case BotStatusRunning:
		return "running"
	case BotStatusStopping:
		return "stopping"
	case BotStatusStopped:
		return "stopped"
	case BotStatusRestarting:
		return "restarting"
	case BotStatusDeleting:
		return "deleting"
	case BotStatusError:
		return "error"
	case BotStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

func (s Status) Int32() int32 {
	return int32(s)
}

type Name string

func NewName(name string) (Name, error) {
	trimmed := strings.TrimSpace(name)
	lenTrimmed := len(trimmed)

	if lenTrimmed < 3 {
		return "", ErrBotNameMustBeLonger
	}
	if lenTrimmed > 32 {
		return "", ErrBotNameMustBeShorter
	}

	return Name(trimmed), nil
}

func (n Name) String() string {
	return string(n)
}

type ApiKey string

func NewApiKey(apiKey string) (ApiKey, error) {
	trimmed := strings.TrimSpace(apiKey)

	if trimmed == "" {
		return "", ErrInvalidApiKey
	}

	return ApiKey(trimmed), nil
}

func (ak ApiKey) String() string {
	return "[REDACTED]"
}

func (ak ApiKey) Raw() string {
	return string(ak)
}

package botdomain

import "errors"

var (
	ErrInvalidBotStatus = errors.New("invalid bot_status")
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

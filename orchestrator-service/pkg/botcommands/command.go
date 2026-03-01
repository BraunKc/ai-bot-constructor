package botcommands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CommandType string

const (
	CommandCreate  CommandType = "bot.create"
	CommandStart   CommandType = "bot.start"
	CommandStop    CommandType = "bot.stop"
	CommandRestart CommandType = "bot.restart"
	CommandDelete  CommandType = "bot.delete"
)

type Command struct {
	ID        uuid.UUID   `json:"id"`
	BotID     uuid.UUID   `json:"bot_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Type      CommandType `json:"type"`
	Payload   []byte      `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

type CreatePayload struct {
	Name   string `json:"name"`
	ApiKey string `json:"api_key"` // SECURITY: must encrypt in prod
}

func NewCommand(botID, userID uuid.UUID, cmdType CommandType, payload any) (*Command, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &Command{
		ID:        uuid.New(),
		BotID:     botID,
		UserID:    userID,
		Type:      cmdType,
		Payload:   payloadBytes,
		Timestamp: time.Now().UTC(),
	}, nil
}

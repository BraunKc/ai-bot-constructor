package botdomain

import (
	"errors"

	"github.com/google/uuid"
)

// SECURITY: API key must be encrypted in prod
type Bot struct {
	id     uuid.UUID
	userID uuid.UUID
	status Status
	name   Name
	apiKey ApiKey
}

var (
	ErrInvalidStorageData = errors.New("invalid storage data")
)

func NewBot(userID uuid.UUID, name string, apiKey string) (*Bot, error) {
	n, err := NewName(name)
	if err != nil {
		return nil, err
	}
	ak, err := NewApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		id:     uuid.New(),
		userID: userID,
		status: BotStatusCreating,
		name:   n,
		apiKey: ak,
	}, nil
}

// USE ONLY FOR CREATING BOT FROM REPOSITORY!!!
func RestoreBot(id, userID uuid.UUID, status int32, name string, apiKey string) (*Bot, error) {
	s, err := NewStatus(status)
	if err != nil {
		return nil, ErrInvalidStorageData
	}
	n, err := NewName(name)
	if err != nil {
		return nil, ErrInvalidStorageData
	}
	ak, err := NewApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return &Bot{
		id:     id,
		userID: userID,
		status: s,
		name:   n,
		apiKey: ak,
	}, nil
}

func (b *Bot) ID() uuid.UUID {
	return b.id
}

func (b *Bot) UserID() uuid.UUID {
	return b.userID
}

func (b *Bot) Status() Status {
	return b.status
}

func (b *Bot) Name() Name {
	return b.name
}

func (b *Bot) ChangeStatus(status int32) error {
	s, err := NewStatus(status)
	if err != nil {
		return err
	}
	b.status = s

	return nil
}

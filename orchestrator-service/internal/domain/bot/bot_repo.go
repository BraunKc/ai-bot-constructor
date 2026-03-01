package botdomain

import (
	"context"

	"github.com/google/uuid"
)

type BotRepo interface {
	Create(ctx context.Context, bot *Bot) error
	GetByID(ctx context.Context, id uuid.UUID) (*Bot, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Bot, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	Close() error
}

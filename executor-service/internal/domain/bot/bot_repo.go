package botdomain

import (
	"context"

	"github.com/google/uuid"
)

type BotRepo interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	UpdateError(ctx context.Context, id uuid.UUID, error string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Close() error
}

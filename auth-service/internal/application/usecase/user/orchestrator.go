package userusecase

import (
	"context"

	"github.com/google/uuid"
)

type OrchestratorClient interface {
	DeleteAllBots(ctx context.Context, userID uuid.UUID) error
	Close() error
}

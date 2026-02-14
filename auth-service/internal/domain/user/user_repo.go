package userdomain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username Username) (*User, error)
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUsername(ctx context.Context, id uuid.UUID, newUsername Username) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

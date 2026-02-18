package userusecase

import (
	"context"

	userdto "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/dto/user"
)

type AuthClient interface {
	Register(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	Login(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	GetUser(ctx context.Context) (*userdto.User, error)
	UpdateUser(ctx context.Context, req *userdto.UpdateUserReq) (*userdto.User, error)
	DeleteUser(ctx context.Context) error
}

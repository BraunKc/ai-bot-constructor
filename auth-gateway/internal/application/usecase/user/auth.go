package userusecase

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	userdto "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/dto/user"
	autherrors "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/infra/grpc/auth/errors"
)

type AuthClient interface {
	Register(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	Login(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	GetUser(ctx context.Context) (*userdto.User, error)
	UpdateUser(ctx context.Context, req *userdto.UpdateUserReq) (*userdto.User, error)
	DeleteUser(ctx context.Context) error
	Close() error
}

type UserUsecase interface {
	Register(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	Login(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error)
	GetUser(ctx context.Context) (*userdto.User, error)
	UpdateUser(ctx context.Context, req *userdto.UpdateUserReq) (*userdto.User, error)
	DeleteUser(ctx context.Context) error
}

type userUsecase struct {
	authClient AuthClient
	log        *slog.Logger
}

func New(authClient AuthClient, log *slog.Logger) UserUsecase {
	return &userUsecase{
		authClient: authClient,
		log:        log,
	}
}

func (uu *userUsecase) Register(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error) {
	uu.log.Debug("initiated register", slog.String("username", req.Username))

	resp, err := uu.authClient.Register(ctx, req)
	if err != nil {
		uu.logIfInternal(err, slog.String("method", "register"))

		return nil, err
	}

	return resp, nil
}

func (uu *userUsecase) Login(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error) {
	uu.log.Debug("initiated login", slog.String("username", req.Username))

	resp, err := uu.authClient.Login(ctx, req)
	if err != nil {
		uu.logIfInternal(err, slog.String("method", "login"))

		return nil, err
	}

	return resp, nil
}

func (uu *userUsecase) GetUser(ctx context.Context) (*userdto.User, error) {
	uu.log.Debug("initiated get user")

	resp, err := uu.authClient.GetUser(ctx)
	if err != nil {
		uu.logIfInternal(err, slog.String("method", "get user"))

		return nil, err
	}

	return resp, nil
}

func (uu *userUsecase) UpdateUser(ctx context.Context, req *userdto.UpdateUserReq) (*userdto.User, error) {
	uu.log.Debug("initiated update user", slog.String("new_username", req.NewUsername))

	resp, err := uu.authClient.UpdateUser(ctx, req)
	if err != nil {
		uu.logIfInternal(err, slog.String("method", "update user"))

		return nil, err
	}

	return resp, nil
}

func (uu *userUsecase) DeleteUser(ctx context.Context) error {
	uu.log.Debug("initiated delete user")

	if err := uu.authClient.DeleteUser(ctx); err != nil {
		uu.logIfInternal(err, slog.String("method", "delete user"))

		return err
	}

	return nil
}

func (uu *userUsecase) logIfInternal(err error, attrs ...slog.Attr) {
	if err == nil {
		return
	}

	var appError *autherrors.AppError
	if errors.As(err, &appError) {
		if appError.HTTPStatus == http.StatusInternalServerError {
			args := append([]any{slog.Any("err", err)}, attrs)
			uu.log.Error("internal error", args...)
		}
	}
}

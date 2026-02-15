package userusecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	userdto "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/dto/user"
	userdomain "github.com/braunkc/ai-bot-constructor/auth-service/internal/domain/user"
	hasherinfra "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/hasher"
	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(ctx context.Context, req *userdto.AuthReq) (string, error)
	Login(ctx context.Context, req *userdto.AuthReq) (string, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*userdto.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, req *userdto.UpdateUserReq) (*userdto.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type userUsecase struct {
	userRepo           userdomain.UserRepo
	orchestratorClient OrchestratorClient
	tokenManager       TokenManager
	hasher             *hasherinfra.Hasher
	log                *slog.Logger
}

func New(userRepo userdomain.UserRepo, orchestratorClient OrchestratorClient,
	tokenManager TokenManager, hasher *hasherinfra.Hasher, log *slog.Logger) UserUsecase {
	return &userUsecase{
		userRepo:           userRepo,
		orchestratorClient: orchestratorClient,
		tokenManager:       tokenManager,
		hasher:             hasher,
		log:                log,
	}
}

func (uu *userUsecase) Register(ctx context.Context, req *userdto.AuthReq) (string, error) {
	uu.log.Debug("register user", slog.String("username", req.Username))

	user, err := userdomain.NewUser(req.Username, req.Password, uu.hasher)
	if err != nil {
		return "", err
	}

	if err := uu.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	return uu.tokenManager.GenerateAccessToken(user.ID().String())
}

func (uu *userUsecase) Login(ctx context.Context, req *userdto.AuthReq) (string, error) {
	uu.log.Debug("login user", slog.String("username", req.Username))

	username, err := userdomain.NewUsername(req.Username)
	if err != nil {
		return "", err
	}

	user, err := uu.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, userdomain.ErrRecordNotFound) {
			return "", userdomain.ErrInvalidUsernameOrPassword
		}

		return "", err
	}

	if !user.CheckPassword(req.Password, uu.hasher) {
		return "", userdomain.ErrInvalidUsernameOrPassword
	}

	return uu.tokenManager.GenerateAccessToken(user.ID().String())
}

func (uu *userUsecase) GetUser(ctx context.Context, userID uuid.UUID) (*userdto.User, error) {
	uu.log.Debug("getting user", slog.Any("user_id", userID))

	user, err := uu.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &userdto.User{
		ID:       user.ID().String(),
		Username: user.Username().String(),
	}, nil
}

func (uu *userUsecase) UpdateUser(ctx context.Context, userID uuid.UUID, req *userdto.UpdateUserReq) (*userdto.User, error) {
	uu.log.Debug("updating user",
		slog.Any("user_id", userID),
		slog.String("new_username", req.NewUsername),
	)

	user, err := uu.userRepo.UpdateUsername(ctx, userID, userdomain.Username(req.NewUsername))
	if err != nil {
		return nil, err
	}

	return &userdto.User{
		ID:       user.ID().String(),
		Username: user.Username().String(),
	}, nil
}

func (uu *userUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	uu.log.Debug("deleting user", slog.Any("user_id", userID))

	if err := uu.orchestratorClient.DeleteAllBots(ctx, userID); err != nil {
		uu.log.Error("failed to delete user bots",
			slog.Any("user_id", userID),
			slog.Any("err", err),
		)

		return fmt.Errorf("failed to delete user bots: %w", err)
	}

	if err := uu.userRepo.Delete(ctx, userID); err != nil {
		uu.log.Error("failed to delete user",
			slog.Any("user_id", userID),
			slog.Any("err", err),
		)

		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

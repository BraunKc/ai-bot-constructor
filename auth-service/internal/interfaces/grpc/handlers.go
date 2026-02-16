package grpcserver

import (
	"context"
	"errors"
	"log/slog"

	authpb "github.com/braunkc/ai-bot-constructor/auth-service/api/auth-service/v1"
	appcontext "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/context"
	userdto "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/dto/user"
	userdomain "github.com/braunkc/ai-bot-constructor/auth-service/internal/domain/user"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GRPCServer) Register(ctx context.Context, req *authpb.AuthReq) (*authpb.Token, error) {
	s.log.Debug("received register request", slog.String("username", req.Username))

	token, err := s.userUsecase.Register(ctx, &userdto.AuthReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &authpb.Token{
		Token: token,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *authpb.AuthReq) (*authpb.Token, error) {
	s.log.Debug("received login request", slog.String("username", req.Username))

	token, err := s.userUsecase.Login(ctx, &userdto.AuthReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &authpb.Token{
		Token: token,
	}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, _ *emptypb.Empty) (*authpb.User, error) {
	s.log.Debug("received get user request", slog.Any("user_id", ctx.Value("user_id")))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	user, err := s.userUsecase.GetUser(ctx, userID)
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &authpb.User{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *GRPCServer) UpdateUser(ctx context.Context, req *authpb.UpdateUserReq) (*authpb.User, error) {
	s.log.Debug("received update user request",
		slog.Any("user_id", ctx.Value("user_id")),
		slog.String("new_username", req.NewUsername),
	)

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	user, err := s.userUsecase.UpdateUser(ctx, userID, &userdto.UpdateUserReq{
		NewUsername: req.NewUsername,
	})
	if err != nil {
		return nil, s.grpcError(err)
	}

	return &authpb.User{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	s.log.Debug("received delete user request", slog.Any("user_id", ctx.Value("user_id")))

	userID, err := s.userIDFromContext(ctx)
	if err != nil {
		return nil, s.grpcError(err)
	}

	if err := s.userUsecase.DeleteUser(ctx, userID); err != nil {
		return nil, s.grpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := appcontext.UserIDFromContext(ctx)
	if !ok {
		return uuid.Nil, errors.New("failed to parse user_id from context")
	}

	return uuid.Parse(userID)
}

func (s *GRPCServer) grpcError(err error) error {
	switch {
	case errors.Is(err, userdomain.ErrDuplicatedKey):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, userdomain.ErrRecordNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, userdomain.ErrInvalidStorageData):
		return status.Error(codes.DataLoss, err.Error())
	case errors.Is(err, userdomain.ErrInvalidUsernameOrPassword):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, userdomain.ErrEmptyUsername),
		errors.Is(err, userdomain.ErrUsernameMustBeLonger),
		errors.Is(err, userdomain.ErrEmptyPassword),
		errors.Is(err, userdomain.ErrPasswordMustBeLonger):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, jwt.ErrTokenExpired),
		errors.Is(err, jwt.ErrInvalidToken),
		errors.Is(err, jwt.ErrInvalidClaims),
		errors.Is(err, jwt.ErrUserIDIsRequired):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		s.log.Error("internal error", slog.Any("err", err))
		return status.Error(codes.Internal, "internal error")
	}
}

package authgrpc

import (
	"context"
	"fmt"
	"log/slog"

	authpb "github.com/braunkc/ai-bot-constructor/auth-gateway/api/auth-service/v1"
	"github.com/braunkc/ai-bot-constructor/auth-gateway/config"
	userdto "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/dto/user"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/application/usecase/user"
	authinterceptors "github.com/braunkc/ai-bot-constructor/auth-gateway/internal/infra/grpc/auth/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthClient
	log    *slog.Logger
}

func NewClient(cfg *config.AuthServiceConfig, log *slog.Logger) (userusecase.AuthClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(authinterceptors.UnaryAuthInterceptor()),
	)
	if err != nil {
		return nil, httpError(err)
	}
	client := authpb.NewAuthClient(conn)

	return &authClient{
		conn:   conn,
		client: client,
		log:    log,
	}, nil
}

func (ac *authClient) Close() error {
	return ac.conn.Close()
}

func (ac *authClient) Register(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error) {
	ac.log.Debug("requesting for register", slog.String("username", req.Username))

	resp, err := ac.client.Register(ctx, &authpb.AuthReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, httpError(err)
	}

	return &userdto.Token{
		Token: resp.Token,
	}, nil
}

func (ac *authClient) Login(ctx context.Context, req *userdto.AuthReq) (*userdto.Token, error) {
	ac.log.Debug("requesting for login", slog.String("username", req.Username))

	resp, err := ac.client.Login(ctx, &authpb.AuthReq{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, httpError(err)
	}

	return &userdto.Token{
		Token: resp.Token,
	}, nil
}

func (ac *authClient) GetUser(ctx context.Context) (*userdto.User, error) {
	ac.log.Debug("requesting for get user")

	resp, err := ac.client.GetUser(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, httpError(err)
	}

	return &userdto.User{
		Username: resp.Username,
	}, nil
}

func (ac *authClient) UpdateUser(ctx context.Context, req *userdto.UpdateUserReq) (*userdto.User, error) {
	ac.log.Debug("requesting for update user", slog.String("new_username", req.NewUsername))

	resp, err := ac.client.UpdateUser(ctx, &authpb.UpdateUserReq{
		NewUsername: req.NewUsername,
	})
	if err != nil {
		return nil, httpError(err)
	}

	return &userdto.User{
		Username: resp.Username,
	}, nil
}

func (ac *authClient) DeleteUser(ctx context.Context) error {
	ac.log.Debug("requesting for delete user")

	_, err := ac.client.DeleteUser(ctx, &emptypb.Empty{})
	if err != nil {
		return httpError(err)
	}

	return nil
}

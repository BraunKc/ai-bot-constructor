package grpcserver

import (
	"log/slog"

	authpb "github.com/braunkc/ai-bot-constructor/auth-service/api/auth-service/v1"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/usecase/user"
	grpcinterceptors "github.com/braunkc/ai-bot-constructor/auth-service/internal/interfaces/grpc/interceptors"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	authpb.UnimplementedAuthServer
	userUsecase userusecase.UserUsecase
	log         *slog.Logger
}

func New(authInterceptors *grpcinterceptors.AuthInterceptor, userUsecase userusecase.UserUsecase, log *slog.Logger) *grpc.Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptors.UnaryServerInterceptor()))
	authpb.RegisterAuthServer(grpcServer, &GRPCServer{
		userUsecase: userUsecase,
		log:         log,
	})

	return grpcServer
}

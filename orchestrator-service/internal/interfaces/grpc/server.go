package grpcserver

import (
	"log/slog"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/orchestrator-service/api/orchestrator-service/v1"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/application/usecase/bot"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/auth"
	grpcinterceptors "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/interfaces/grpc/interceptors"
	"google.golang.org/grpc"
)

type OrchestratorServiceServer struct {
	orchestratorpb.UnimplementedOrchestratorServer
	botUsecase botusecase.BotUsecase
	log        *slog.Logger
}

func New(botUsecase botusecase.BotUsecase, tokenManager *auth.TokenManager, log *slog.Logger) *grpc.Server {
	authInterceptor := grpcinterceptors.AuthInterceptor(map[string]struct{}{}, tokenManager)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)
	orchestratorpb.RegisterOrchestratorServer(grpcServer, &OrchestratorServiceServer{
		botUsecase: botUsecase,
		log:        log,
	})

	return grpcServer
}

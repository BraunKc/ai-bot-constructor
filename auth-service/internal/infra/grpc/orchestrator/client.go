package orchestratorgrpc

import (
	"context"
	"fmt"
	"log/slog"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/auth-service/api/orchestrator-service/v1"
	"github.com/braunkc/ai-bot-constructor/auth-service/config"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/usecase/user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orchestratorClient struct {
	conn   *grpc.ClientConn
	client orchestratorpb.OrchestratorClient
	log    *slog.Logger
}

func NewClient(cfg *config.OrchestratorServiceConfig, log *slog.Logger) (userusecase.OrchestratorClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := orchestratorpb.NewOrchestratorClient(conn)

	return &orchestratorClient{
		client: client,
		log:    log,
	}, nil
}

func (oc *orchestratorClient) Close() error {
	return oc.conn.Close()
}

func (oc *orchestratorClient) DeleteAllBots(ctx context.Context, userID uuid.UUID) error {
	oc.log.Debug("deleting all user bots", slog.Any("user_id", userID))

	resp, err := oc.client.DeleteAllBots(ctx, &orchestratorpb.DeleteAllBotsReq{
		UserId: userID.String(),
	})
	if err != nil {
		return err
	}
	if !resp.AllSucceeded {
		oc.log.Warn("not all user bots were deleted", slog.Any("user_id", userID))
	}

	return nil
}

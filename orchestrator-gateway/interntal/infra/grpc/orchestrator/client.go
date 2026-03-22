package orchestratorgrpc

import (
	"context"
	"fmt"
	"log/slog"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/api/orchestrator-service/v1"
	"github.com/braunkc/ai-bot-constructor/orchestrator-gateway/config"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/usecase/bot"
	orchestratorerrors "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/infra/grpc/orchestrator/errors"
	orchestratorinterceptors "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/infra/grpc/orchestrator/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orchestratorClient struct {
	conn   *grpc.ClientConn
	client orchestratorpb.OrchestratorClient
	log    *slog.Logger
}

func NewClient(cfg *config.OrchestratorServiceConfig, log *slog.Logger) (botusecase.OrchestratorClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(orchestratorinterceptors.UnaryAuthInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create conn to orchestrator service")
	}

	return &orchestratorClient{
		conn:   conn,
		client: orchestratorpb.NewOrchestratorClient(conn),
		log:    log,
	}, nil
}

func (oc *orchestratorClient) CreateBot(ctx context.Context, name, systemPrompt, apiKey string) (*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for create bot", slog.String("name", name))

	resp, err := oc.client.CreateBot(ctx, &orchestratorpb.CreateBotReq{
		Name:         name,
		SystemPrompt: systemPrompt,
		ApiKey:       apiKey,
	})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp, nil
}

func (oc *orchestratorClient) GetBot(ctx context.Context, id string) (*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for get bot", slog.String("bot_id", id))

	resp, err := oc.client.GetBot(ctx, &orchestratorpb.GetBotReq{Id: id})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp, nil
}

func (oc *orchestratorClient) GetAllBots(ctx context.Context) ([]*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for get all bots")

	resp, err := oc.client.GetAllBots(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp.Bots, nil
}

func (oc *orchestratorClient) StartBot(ctx context.Context, id string) (*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for start bot", slog.String("bot_id", id))

	resp, err := oc.client.StartBot(ctx, &orchestratorpb.StartBotReq{Id: id})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp, nil
}

func (oc *orchestratorClient) StartBots(ctx context.Context, ids []string) ([]*orchestratorpb.Bot, bool, error) {
	oc.log.Debug("requesting for start bots", slog.Any("botds_ids", ids))

	resp, err := oc.client.StartBots(ctx, &orchestratorpb.StartBotsReq{Ids: ids})
	if err != nil {
		return nil, false, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp.Bots, resp.AllSucceeded, nil
}

func (oc *orchestratorClient) StopBot(ctx context.Context, id string) (*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for stop bot", slog.String("bot_id", id))

	resp, err := oc.client.StopBot(ctx, &orchestratorpb.StopBotReq{Id: id})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp, nil
}

func (oc *orchestratorClient) StopBots(ctx context.Context, ids []string) ([]*orchestratorpb.Bot, bool, error) {
	oc.log.Debug("requesting for stop bots", slog.Any("bots_ids", ids))

	resp, err := oc.client.StopBots(ctx, &orchestratorpb.StopBotsReq{Ids: ids})
	if err != nil {
		return nil, false, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp.Bots, resp.AllSucceeded, nil
}

func (oc *orchestratorClient) RestartBot(ctx context.Context, id string) (*orchestratorpb.Bot, error) {
	oc.log.Debug("requesting for restart bot", slog.String("bot_id", id))

	resp, err := oc.client.RestartBot(ctx, &orchestratorpb.RestartBotReq{Id: id})
	if err != nil {
		return nil, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp, nil
}

func (oc *orchestratorClient) DeleteBot(ctx context.Context, id string) error {
	oc.log.Debug("requesting for delete bot", slog.String("bot_id", id))

	_, err := oc.client.DeleteBot(ctx, &orchestratorpb.DeleteBotReq{Id: id})
	if err != nil {
		return orchestratorerrors.GRPCToHTTPError(err)
	}

	return nil
}

func (oc *orchestratorClient) DeleteBots(ctx context.Context, ids []string) (bool, error) {
	oc.log.Debug("requesting for delete bots", slog.Any("bots_ids", ids))

	resp, err := oc.client.DeleteBots(ctx, &orchestratorpb.DeleteBotsReq{Ids: ids})
	if err != nil {
		return false, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp.AllSucceeded, nil
}

func (oc *orchestratorClient) DeleteAllBots(ctx context.Context) (bool, error) {
	oc.log.Debug("requesting for delete all bots")

	resp, err := oc.client.DeleteAllBots(ctx, &emptypb.Empty{})
	if err != nil {
		return false, orchestratorerrors.GRPCToHTTPError(err)
	}

	return resp.AllSucceeded, nil
}

func (oc *orchestratorClient) Close() error {
	if oc.conn != nil {
		return oc.conn.Close()
	}
	return nil
}

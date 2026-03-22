package botusecase

import (
	"context"
	"log/slog"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/api/orchestrator-service/v1"
	botdto "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/dto/bot"
)

type OrchestratorClient interface {
	CreateBot(ctx context.Context, name, systemPrompt, apiKey string) (*orchestratorpb.Bot, error)
	GetBot(ctx context.Context, id string) (*orchestratorpb.Bot, error)
	GetAllBots(ctx context.Context) ([]*orchestratorpb.Bot, error)
	StartBot(ctx context.Context, id string) (*orchestratorpb.Bot, error)
	StartBots(ctx context.Context, ids []string) ([]*orchestratorpb.Bot, bool, error)
	StopBot(ctx context.Context, id string) (*orchestratorpb.Bot, error)
	StopBots(ctx context.Context, ids []string) ([]*orchestratorpb.Bot, bool, error)
	RestartBot(ctx context.Context, id string) (*orchestratorpb.Bot, error)
	DeleteBot(ctx context.Context, id string) error
	DeleteBots(ctx context.Context, ids []string) (bool, error)
	DeleteAllBots(ctx context.Context) (bool, error)
	Close() error
}

type BotUsecase interface {
	CreateBot(ctx context.Context, req *botdto.CreateBotRequest) (*botdto.Bot, error)
	GetBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error)
	GetAllBots(ctx context.Context) ([]*botdto.Bot, error)
	StartBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error)
	StartBots(ctx context.Context, req *botdto.IDsRequest) (*botdto.OperationResponse, error)
	StopBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error)
	StopBots(ctx context.Context, req *botdto.IDsRequest) (*botdto.OperationResponse, error)
	RestartBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error)
	DeleteBot(ctx context.Context, req *botdto.GetBotRequest) error
	DeleteBots(ctx context.Context, req *botdto.IDsRequest) (bool, error)
	DeleteAllBots(ctx context.Context) (bool, error)
}

type botUsecase struct {
	client OrchestratorClient
	log    *slog.Logger
}

func NewBotUsecase(orchestratorClient OrchestratorClient, log *slog.Logger) BotUsecase {
	return &botUsecase{
		client: orchestratorClient,
		log:    log,
	}
}

func (bu *botUsecase) CreateBot(ctx context.Context, req *botdto.CreateBotRequest) (*botdto.Bot, error) {
	bu.log.Debug("initiated creating bot", slog.String("name", req.Name))

	pbBot, err := bu.client.CreateBot(ctx, req.Name, req.SystemPrompt, req.APIKey)
	if err != nil {
		return nil, err
	}

	return botdto.FromProto(pbBot), nil
}

func (bu *botUsecase) GetBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error) {
	bu.log.Info("initiated getting bot", slog.String("id", req.ID))

	pbBot, err := bu.client.GetBot(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return botdto.FromProto(pbBot), nil
}

func (bu *botUsecase) GetAllBots(ctx context.Context) ([]*botdto.Bot, error) {
	bu.log.Info("initiated getting all bots")

	pbBots, err := bu.client.GetAllBots(ctx)
	if err != nil {
		return nil, err
	}

	return botdto.FromProtoList(pbBots), nil
}

func (bu *botUsecase) StartBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error) {
	bu.log.Info("initiated starting bot", slog.String("id", req.ID))

	pbBot, err := bu.client.StartBot(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return botdto.FromProto(pbBot), nil
}

func (bu *botUsecase) StartBots(ctx context.Context, req *botdto.IDsRequest) (*botdto.OperationResponse, error) {
	bu.log.Info("initiated starting bots", slog.Any("ids", req.IDs))

	pbBots, allSucceeded, err := bu.client.StartBots(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &botdto.OperationResponse{
		AllSucceeded: allSucceeded,
		Bots:         botdto.FromProtoList(pbBots),
	}, nil
}

func (bu *botUsecase) StopBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error) {
	bu.log.Info("initiated stopping bot", slog.String("id", req.ID))

	pbBot, err := bu.client.StopBot(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return botdto.FromProto(pbBot), nil
}

func (bu *botUsecase) StopBots(ctx context.Context, req *botdto.IDsRequest) (*botdto.OperationResponse, error) {
	bu.log.Info("initiated stopping bots", slog.Any("ids", req.IDs))

	pbBots, allSucceeded, err := bu.client.StopBots(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &botdto.OperationResponse{
		AllSucceeded: allSucceeded,
		Bots:         botdto.FromProtoList(pbBots),
	}, nil
}

func (bu *botUsecase) RestartBot(ctx context.Context, req *botdto.GetBotRequest) (*botdto.Bot, error) {
	bu.log.Info("initiated restarting bot", slog.String("id", req.ID))

	pbBot, err := bu.client.RestartBot(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return botdto.FromProto(pbBot), nil
}

func (bu *botUsecase) DeleteBot(ctx context.Context, req *botdto.GetBotRequest) error {
	bu.log.Info("initiated deleting bot", slog.String("id", req.ID))

	return bu.client.DeleteBot(ctx, req.ID)
}

func (bu *botUsecase) DeleteBots(ctx context.Context, req *botdto.IDsRequest) (bool, error) {
	bu.log.Info("initiated deleting bots", slog.Any("ids", req.IDs))

	return bu.client.DeleteBots(ctx, req.IDs)
}

func (bu *botUsecase) DeleteAllBots(ctx context.Context) (bool, error) {
	bu.log.Info("initiated deleting all bots")
	return bu.client.DeleteAllBots(ctx)
}

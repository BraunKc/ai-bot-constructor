package botusecase

import (
	"context"
	"log/slog"

	botdto "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/application/dto/bot"
	botdomain "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/domain/bot"
	kafkaproducer "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/repo/kafka/bot/producer"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/pkg/botcommands"
	"github.com/google/uuid"
)

type BotUsecase interface {
	CreateBot(ctx context.Context, userID uuid.UUID, req *botdto.CreateBotReq) (*botdto.Bot, error)
	GetBot(ctx context.Context, userID uuid.UUID, req *botdto.GetBotReq) (*botdto.Bot, error)
	GetAllBots(ctx context.Context, userID uuid.UUID) ([]*botdto.Bot, error)
	StopBot(ctx context.Context, userID uuid.UUID, req *botdto.StopBotReq) (*botdto.Bot, error)
	StopBots(ctx context.Context, userID uuid.UUID, req *botdto.StopBotsReq) (*botdto.StopBotsResp, error)
	StartBot(ctx context.Context, userID uuid.UUID, req *botdto.StartBotReq) (*botdto.Bot, error)
	RestartBot(ctx context.Context, userID uuid.UUID, req *botdto.RestartBotReq) (*botdto.Bot, error)
	DeleteBot(ctx context.Context, userID uuid.UUID, req *botdto.DeleteBotReq) error
	DeleteBots(ctx context.Context, userID uuid.UUID, req *botdto.DeleteBotsReq) (*botdto.DeleteBotsResp, error)
	DeleteAllBots(ctx context.Context, userID uuid.UUID) (*botdto.DeleteAllBotsResp, error)
}

type botUsecase struct {
	botRepo       botdomain.BotRepo
	kafkaProducer kafkaproducer.KafkaProducer
	log           *slog.Logger
}

func New(botRepo botdomain.BotRepo, kafkaProducer kafkaproducer.KafkaProducer, log *slog.Logger) BotUsecase {
	return &botUsecase{
		botRepo:       botRepo,
		kafkaProducer: kafkaProducer,
		log:           log,
	}
}

func (bu *botUsecase) CreateBot(ctx context.Context, userID uuid.UUID, req *botdto.CreateBotReq) (*botdto.Bot, error) {
	bu.log.Debug("creating bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_name", req.Name),
	)

	bot, err := botdomain.NewBot(userID, req.Name, req.ApiKey)
	if err != nil {
		return nil, err
	}

	if err := bu.botRepo.Create(ctx, bot); err != nil {
		return nil, err
	}

	cmd, err := botcommands.NewCommand(bot.UserID(), userID, botcommands.CommandCreate, botcommands.CreatePayload{Name: bot.Name().String(), ApiKey: bot.ApiKey().Raw()})
	if err != nil {
		return nil, err
	}

	if err := bu.kafkaProducer.Produce(ctx, cmd); err != nil {
		return nil, err
	}

	return bu.DomainToDTOModel(bot), nil
}

func (bu *botUsecase) GetBot(ctx context.Context, userID uuid.UUID, req *botdto.GetBotReq) (*botdto.Bot, error) {
	bu.log.Debug("getting bot by id",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", req.ID.String()),
	)

	bot, err := bu.botRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if bot.UserID() != userID {
		return nil, botdomain.ErrNotEnoughRights
	}

	return bu.DomainToDTOModel(bot), nil
}

func (bu *botUsecase) GetAllBots(ctx context.Context, userID uuid.UUID) ([]*botdto.Bot, error) {
	bu.log.Debug("getting all bots by user_id", slog.String("user_id", userID.String()))

	bots, err := bu.botRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	botsDTO := make([]*botdto.Bot, 0, len(bots))
	for _, bot := range bots {
		botsDTO = append(botsDTO, bu.DomainToDTOModel(bot))
	}

	return botsDTO, nil
}

func (bu *botUsecase) StopBot(ctx context.Context, userID uuid.UUID, req *botdto.StopBotReq) (*botdto.Bot, error) {
	bu.log.Debug("stopping bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", req.ID.String()),
	)

	bot, err := bu.botRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if bot.UserID() != userID {
		return nil, botdomain.ErrNotEnoughRights
	}

	if err := bot.ChangeStatus(botdomain.BotStatusStopping.Int32()); err != nil {
		return nil, err
	}

	if err := bu.botRepo.UpdateStatus(ctx, bot.ID(), bot.Status()); err != nil {
		return nil, err
	}

	cmd, err := botcommands.NewCommand(req.ID, userID, botcommands.CommandStop, nil)
	if err != nil {
		return nil, err
	}

	return bu.DomainToDTOModel(bot), bu.kafkaProducer.Produce(ctx, cmd)
}

func (bu *botUsecase) StopBots(ctx context.Context, userID uuid.UUID, req *botdto.StopBotsReq) (*botdto.StopBotsResp, error) {
	bu.log.Debug("stopping bots by id",
		slog.String("user_id", userID.String()),
		slog.Any("bot_ids", req.IDs),
	)

	bots := make([]botdto.Bot, 0, len(req.IDs))
	for _, id := range req.IDs {
		bot, err := bu.botRepo.GetByID(ctx, id)
		if err != nil {
			bu.log.Warn("failed to get bot by id",
				slog.String("bot_id", id.String()),
				slog.Any("err", err),
			)

			continue
		}

		if bot.UserID() != userID {
			continue
		}

		if err := bot.ChangeStatus(botdomain.BotStatusStopping.Int32()); err != nil {
			bu.log.Error("failed to change bot status at domain", slog.Any("err", err))
			continue
		}

		if err := bu.botRepo.UpdateStatus(ctx, bot.ID(), bot.Status()); err != nil {
			bu.log.Error("failed to change bot status at db", slog.Any("err", err))
			continue
		}

		cmd, err := botcommands.NewCommand(bot.ID(), bot.UserID(), botcommands.CommandStop, nil)
		if err != nil {
			bu.log.Error("failed to create stop command", slog.Any("err", err))
			continue
		}

		if err := bu.kafkaProducer.Produce(ctx, cmd); err != nil {
			bu.log.Error("failed to produce stop command", slog.Any("err", err))
			continue
		}

		bots = append(bots, *bu.DomainToDTOModel(bot))
	}

	var allSucceeded bool
	if len(bots) == len(req.IDs) {
		allSucceeded = true
	}

	return &botdto.StopBotsResp{
		Bots:         bots,
		AllSucceeded: allSucceeded,
	}, nil
}

func (bu *botUsecase) StartBot(ctx context.Context, userID uuid.UUID, req *botdto.StartBotReq) (*botdto.Bot, error) {
	bu.log.Debug("starting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", req.ID.String()),
	)

	bot, err := bu.botRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if bot.UserID() != userID {
		return nil, botdomain.ErrNotEnoughRights
	}

	if err := bot.ChangeStatus(botdomain.BotStatusStarting.Int32()); err != nil {
		return nil, err
	}

	if err := bu.botRepo.UpdateStatus(ctx, bot.ID(), bot.Status()); err != nil {
		return nil, err
	}

	cmd, err := botcommands.NewCommand(bot.ID(), userID, botcommands.CommandStart, nil)
	if err != nil {
		return nil, err
	}

	return bu.DomainToDTOModel(bot), bu.kafkaProducer.Produce(ctx, cmd)
}

func (bu *botUsecase) RestartBot(ctx context.Context, userID uuid.UUID, req *botdto.RestartBotReq) (*botdto.Bot, error) {
	bu.log.Debug("restarting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", req.ID.String()),
	)

	bot, err := bu.botRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if bot.UserID() != userID {
		return nil, botdomain.ErrNotEnoughRights
	}

	if err := bot.ChangeStatus(botdomain.BotStatusRestarting.Int32()); err != nil {
		return nil, err
	}

	if err := bu.botRepo.UpdateStatus(ctx, bot.ID(), bot.Status()); err != nil {
		return nil, err
	}

	cmd, err := botcommands.NewCommand(bot.ID(), userID, botcommands.CommandRestart, nil)
	if err != nil {
		return nil, err
	}

	return bu.DomainToDTOModel(bot), bu.kafkaProducer.Produce(ctx, cmd)
}

func (bu *botUsecase) DeleteBot(ctx context.Context, userID uuid.UUID, req *botdto.DeleteBotReq) error {
	if _, err := bu.deleteBotFullCycle(ctx, userID, req.ID); err != nil {
		return err
	}

	return nil
}

func (bu *botUsecase) DeleteBots(ctx context.Context, userID uuid.UUID, req *botdto.DeleteBotsReq) (*botdto.DeleteBotsResp, error) {
	bu.log.Debug("deleting bots",
		slog.String("user_id", userID.String()),
		slog.Any("bot_ids", req.IDs),
	)

	successfulIDs := make([]uuid.UUID, 0, len(req.IDs))
	for _, id := range req.IDs {
		bot, err := bu.deleteBotFullCycle(ctx, userID, id)
		if err != nil {
			continue
		}

		successfulIDs = append(successfulIDs, bot.ID())
	}

	return &botdto.DeleteBotsResp{
		AllSucceeded: len(req.IDs) == len(successfulIDs),
	}, nil
}

func (bu *botUsecase) DeleteAllBots(ctx context.Context, userID uuid.UUID) (*botdto.DeleteAllBotsResp, error) {
	bu.log.Debug("deleting all user bots", slog.String("user_id", userID.String()))

	bots, err := bu.botRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	successfulIDs := make([]uuid.UUID, 0, len(bots))
	for _, bot := range bots {
		bot, err := bu.deleteBotFullCycle(ctx, userID, bot.ID())
		if err != nil {
			continue
		}

		successfulIDs = append(successfulIDs, bot.ID())
	}

	return &botdto.DeleteAllBotsResp{
		AllSucceeded: len(bots) == len(successfulIDs),
	}, nil
}

func (bu *botUsecase) deleteBotFullCycle(ctx context.Context, userID, botID uuid.UUID) (*botdomain.Bot, error) {
	bu.log.Debug("deleting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
	)

	bot, err := bu.botRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, err
	}

	if bot.UserID() != userID {
		return nil, botdomain.ErrNotEnoughRights
	}

	if err := bot.ChangeStatus(botdomain.BotStatusDeleting.Int32()); err != nil {
		bu.log.Error("failed to change bot status", slog.Any("err", err))
		return nil, err
	}

	if err := bu.botRepo.UpdateStatus(ctx, bot.ID(), bot.Status()); err != nil {
		bu.log.Error("failed to update bot status", slog.Any("err", err))
		return nil, err
	}

	cmd, err := botcommands.NewCommand(bot.ID(), userID, botcommands.CommandDelete, nil)
	if err != nil {
		bu.log.Error("failed to create command", slog.Any("err", err))
		return nil, err
	}

	if err := bu.kafkaProducer.Produce(ctx, cmd); err != nil {
		bu.log.Error("failed to produce message", slog.Any("err", err))
		return nil, err
	}

	return bot, nil
}

func (bu *botUsecase) DomainToDTOModel(bot *botdomain.Bot) *botdto.Bot {
	return &botdto.Bot{
		ID:        bot.ID(),
		UserID:    bot.UserID(),
		BotStatus: bot.Status().Int32(),
		Name:      bot.Name().String(),
		LastError: bot.LastError(),
		CreatedAt: bot.CreatedAt(),
		UpdatedAt: bot.UpdatedAt(),
	}
}

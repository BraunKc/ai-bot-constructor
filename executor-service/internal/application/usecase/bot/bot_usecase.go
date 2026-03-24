package botusecase

import (
	"context"
	"fmt"
	"log/slog"

	botdomain "github.com/braunkc/ai-bot-constructor/executor-service/internal/domain/bot"
	"github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/docker"
	"github.com/google/uuid"
)

type BotUsecase interface {
	Create(ctx context.Context, userID, botID uuid.UUID, botName, tgBotToken, systemPrompt string) error
	Start(ctx context.Context, userID, botID uuid.UUID) error
	Stop(ctx context.Context, userID, botID uuid.UUID) error
	Restart(ctx context.Context, userID, botID uuid.UUID) error
	Delete(ctx context.Context, userID, botID uuid.UUID) error
}

type botUsecase struct {
	openRouterToken string
	docker          docker.Docker
	botRepo         botdomain.BotRepo
	log             *slog.Logger
}

func New(openRouterToken string, docker docker.Docker, botRepo botdomain.BotRepo, log *slog.Logger) BotUsecase {
	return &botUsecase{
		openRouterToken: openRouterToken,
		docker:          docker,
		botRepo:         botRepo,
		log:             log,
	}
}

func (bu *botUsecase) Create(ctx context.Context, userID, botID uuid.UUID, botName, tgBotToken, systemPrompt string) error {
	bu.log.Debug("creating bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("bot_name", botName),
	)

	if err := bu.docker.CreateContainer(ctx, userID, botID, botName, tgBotToken, bu.openRouterToken, systemPrompt); err != nil {
		if err := bu.botRepo.UpdateError(ctx, botID, err.Error()); err != nil {
			return fmt.Errorf("failed to update bot error: %w", err)
		}

		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := bu.botRepo.UpdateStatus(ctx, botID, botdomain.BotStatusStopped); err != nil {
		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (bu *botUsecase) Start(ctx context.Context, userID, botID uuid.UUID) error {
	bu.log.Debug("starting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
	)

	if err := bu.docker.StartContainer(ctx, userID, botID); err != nil {
		if err := bu.botRepo.UpdateError(ctx, botID, err.Error()); err != nil {
			return fmt.Errorf("failed to update bot error: %w", err)
		}

		return fmt.Errorf("failed to start container: %w", err)
	}

	if err := bu.botRepo.UpdateStatus(ctx, botID, botdomain.BotStatusRunning); err != nil {
		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (bu *botUsecase) Stop(ctx context.Context, userID, botID uuid.UUID) error {
	bu.log.Debug("stopping bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
	)

	if err := bu.docker.StopContainer(ctx, userID, botID); err != nil {
		if err := bu.botRepo.UpdateError(ctx, botID, err.Error()); err != nil {
			return fmt.Errorf("failed to update bot error: %w", err)
		}

		return fmt.Errorf("failed to stop container: %w", err)
	}

	if err := bu.botRepo.UpdateStatus(ctx, botID, botdomain.BotStatusStopped); err != nil {
		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (bu *botUsecase) Restart(ctx context.Context, userID, botID uuid.UUID) error {
	bu.log.Debug("restarting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
	)

	if err := bu.docker.RestartContainer(ctx, userID, botID); err != nil {
		if err := bu.botRepo.UpdateError(ctx, botID, err.Error()); err != nil {
			return fmt.Errorf("failed to update bot error: %w", err)
		}

		return fmt.Errorf("failed to restart container: %w", err)
	}

	if err := bu.botRepo.UpdateStatus(ctx, botID, botdomain.BotStatusRunning); err != nil {
		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (bu *botUsecase) Delete(ctx context.Context, userID, botID uuid.UUID) error {
	bu.log.Debug("deleting bot",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
	)

	if err := bu.docker.DeleteContainer(ctx, userID, botID); err != nil {
		if err := bu.botRepo.UpdateError(ctx, botID, err.Error()); err != nil {
			return fmt.Errorf("failed to update bot error: %w", err)
		}

		return fmt.Errorf("failed to delete container: %w", err)
	}

	if err := bu.botRepo.Delete(ctx, botID); err != nil {
		return fmt.Errorf("failed to delete bot: %w", err)
	}

	return nil
}

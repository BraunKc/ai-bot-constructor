package kafkaconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/braunkc/ai-bot-constructor/executor-service/config"
	botusecase "github.com/braunkc/ai-bot-constructor/executor-service/internal/application/usecase/bot"
	botcommands "github.com/braunkc/ai-bot-constructor/executor-service/pkg/botcommands"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer interface {
	Consume(ctx context.Context)
	Close() error
}

type kafkaConsumer struct {
	botUsecase botusecase.BotUsecase
	reader     *kafka.Reader
	log        *slog.Logger
}

func New(botUsecase botusecase.BotUsecase, cfg *config.KafkaConfig, log *slog.Logger) KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Topic:   "bots",
		GroupID: "executor-service",
	})

	return &kafkaConsumer{
		botUsecase: botUsecase,
		reader:     reader,
		log:        log,
	}
}

func (kc *kafkaConsumer) Consume(ctx context.Context) {
	for {
		msg, err := kc.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			kc.log.Error("failed to fetch message", slog.Any("err", err))
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		go func() {
			<-ctx.Done()
			cancel()
		}()

		kc.log.Debug("received message", slog.String("topic", msg.Topic), slog.Int64("offset", msg.Offset))

		var cmd botcommands.Command
		if err := json.Unmarshal(msg.Value, &cmd); err != nil {
			kc.log.Error("failed to unmarshal command", slog.Any("err", err))
			if err := kc.reader.CommitMessages(ctx, msg); err != nil {
				kc.log.Error("failed to commit message after unmarshal error", slog.Any("err", err))
			}

			continue
		}

		if err := kc.processCommand(ctx, cmd); err != nil {
			kc.log.Error("failed to process command", slog.Any("err", err),
				slog.String("command_id", cmd.ID.String()),
				slog.String("type", string(cmd.Type)),
			)
		}

		if err := kc.reader.CommitMessages(ctx, msg); err != nil {
			kc.log.Error("failed to commit message", slog.Any("err", err))
		}
	}
}

func (kc *kafkaConsumer) processCommand(ctx context.Context, cmd botcommands.Command) error {
	switch cmd.Type {
	case botcommands.CommandCreate:
		var payload botcommands.CreatePayload
		if err := json.Unmarshal(cmd.Payload, &payload); err != nil {
			return fmt.Errorf("failed to unmarshal create payload: %w", err)
		}
		return kc.botUsecase.Create(ctx, cmd.UserID, cmd.BotID, payload.Name, payload.ApiKey, payload.SystemPrompt)

	case botcommands.CommandStart:
		return kc.botUsecase.Start(ctx, cmd.UserID, cmd.BotID)

	case botcommands.CommandStop:
		return kc.botUsecase.Stop(ctx, cmd.UserID, cmd.BotID)

	case botcommands.CommandRestart:
		return kc.botUsecase.Restart(ctx, cmd.UserID, cmd.BotID)

	case botcommands.CommandDelete:
		return kc.botUsecase.Delete(ctx, cmd.UserID, cmd.BotID)

	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

func (kc *kafkaConsumer) Close() error {
	return kc.reader.Close()
}

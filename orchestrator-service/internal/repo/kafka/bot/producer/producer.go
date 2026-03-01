package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/braunkc/ai-bot-constructor/orchestrator-service/config"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/pkg/botcommands"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(ctx context.Context, cmd *botcommands.Command) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func New(cfg *config.KafkaConfig, log *slog.Logger) KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Topic:    "bots",
		Balancer: &kafka.Hash{},
		Logger:   kafka.LoggerFunc(log.Debug),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Error("kafka writer error", append([]interface{}{"msg", msg}, args...)...)
		}),
	})

	return &kafkaProducer{
		writer: writer,
		log:    log,
	}
}

func (kp *kafkaProducer) Produce(ctx context.Context, cmd *botcommands.Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal cmd: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(cmd.BotID.String()),
		Value: data,
		Time:  cmd.Timestamp,
	}

	if err := kp.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	kp.log.Debug("command produced",
		slog.Any("command_id", cmd.ID),
		slog.Any("command_typee", cmd.Type),
		slog.Any("bot_id", cmd.BotID),
	)

	return nil
}

func (kp *kafkaProducer) Close() error {
	return kp.writer.Close()
}

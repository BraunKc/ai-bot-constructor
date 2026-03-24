package main

import (
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/executor-service/config"
	botusecase "github.com/braunkc/ai-bot-constructor/executor-service/internal/application/usecase/bot"
	"github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/docker"
	kafkaconsumer "github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/kafka/bot/consumer"
	"github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/repo/database"
	botpersistence "github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/repo/persistence"
	"github.com/braunkc/ai-bot-constructor/executor-service/pkg/log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("failed to load .env file", slog.Any("err", err))
		os.Exit(1)
	}

	cfg, err := config.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		slog.Error("failed to create config", slog.Any("err", err))
		os.Exit(1)
	}

	logCfg := log.Config{
		Service:    cfg.Logger.Service,
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}
	h, err := log.NewHandler(&logCfg)
	if err != nil {
		slog.Error("failed to create log handler", slog.Any("err", err))
		os.Exit(1)
	}
	log := slog.New(h)

	docker, err := docker.New(log)
	if err != nil {
		log.Error("failed to init docker", slog.Any("err", err))
		os.Exit(1)
	}

	// if err := d.CreateContainer(context.Background(), "aboba", "botaboba", "aimal", "8705407659:AAE1TcyNRrRxSqgM5JtBkYnm1fuTp7oHcBo", cfg.OpenRouterToken, "ты - шифратор, ты должен переводить мои сообщения в азбуку морзе"); err != nil {
	// 	log.Error("failed to create container", slog.Any("err", err))
	// 	os.Exit(1)
	// }

	// if err := d.StartContainer(context.Background(), "aboba", "botaboba", "aimal"); err != nil {
	// 	log.Error("failed to start container", slog.Any("err", err))
	// 	os.Exit(1)
	// }

	// if err := d.RestartContainer(context.Background(), "aboba", "botaboba", "aimal"); err != nil {
	// 	log.Error("failed to restart container", slog.Any("err", err))
	// 	os.Exit(1)
	// }

	// if err := d.StopContainer(context.Background(), "aboba", "botaboba", "aimal"); err != nil {
	// 	log.Error("failed to stop container", slog.Any("err", err))
	// 	os.Exit(1)
	// }

	// if err := d.DeleteContainer(context.Background(), "aboba", "botaboba", "aimal"); err != nil {
	// 	log.Error("failed to delete container", slog.Any("err", err))
	// 	os.Exit(1)
	// }

	db, err := database.New(&cfg.DB)
	if err != nil {
		log.Error("failed to init database", slog.Any("err", err))
		os.Exit(1)
	}

	botRepo := botpersistence.NewBotRepo(db, log)

	botUsecase := botusecase.New(cfg.OpenRouterToken, docker, botRepo, log)

	kafkaConsumer := kafkaconsumer.New(botUsecase, &cfg.Kafka, log)

	log.Info("started")
	kafkaConsumer.Consume()
}

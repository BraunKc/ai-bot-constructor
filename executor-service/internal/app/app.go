package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/braunkc/ai-bot-constructor/executor-service/config"
	botusecase "github.com/braunkc/ai-bot-constructor/executor-service/internal/application/usecase/bot"
	botdomain "github.com/braunkc/ai-bot-constructor/executor-service/internal/domain/bot"
	"github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/docker"
	kafkaconsumer "github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/kafka/bot/consumer"
	"github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/repo/database"
	botpersistence "github.com/braunkc/ai-bot-constructor/executor-service/internal/infra/repo/persistence"
	"github.com/braunkc/ai-bot-constructor/executor-service/pkg/log"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	db            *gorm.DB
	docker        docker.Docker
	botRepo       botdomain.BotRepo
	kafkaConsumer kafkaconsumer.KafkaConsumer
	botUsecase    botusecase.BotUsecase
	cfg           *config.Config
	log           *slog.Logger
}

func (app *App) initConfig(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg, err := config.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}
	app.cfg = cfg
	slog.Debug("cfg created")

	return nil
}

func (app *App) initLogger() error {
	var loggerOutputType log.OutputType
	switch app.cfg.Logger.OutputType {
	case config.Console:
		loggerOutputType = log.Console
	case config.File:
		loggerOutputType = log.File
	case config.Both:
		loggerOutputType = log.Both
	}

	var loggerLevel slog.Level
	switch app.cfg.Logger.Level {
	case config.Debug:
		loggerLevel = slog.LevelDebug
	case config.Info:
		loggerLevel = slog.LevelInfo
	case config.Warn:
		loggerLevel = slog.LevelWarn
	case config.Error:
		loggerLevel = slog.LevelError
	}

	logCfg := log.Config{
		Service:    app.cfg.Logger.Service,
		OutputType: loggerOutputType,
		Level:      loggerLevel,
	}
	h, err := log.NewHandler(&logCfg)
	if err != nil {
		return fmt.Errorf("failed to create log handler: %w", err)
	}
	app.log = slog.New(h)

	return nil
}

func (app *App) initDB() error {
	db, err := database.New(&app.cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to create database repository: %w", err)
	}

	app.db = db
	app.log.Info("database inited")

	return nil
}

func (app *App) initDocker() error {
	docker, err := docker.New(app.log)
	if err != nil {
		return fmt.Errorf("failed to init docker: %w", err)
	}

	app.docker = docker
	app.log.Info("docker inited")

	return nil
}

func (app *App) initBotRepo() {
	app.botRepo = botpersistence.NewBotRepo(app.db, app.log)
	app.log.Info("bot repository created")
}

func (app *App) initBotUsecase() {
	app.botUsecase = botusecase.New(app.cfg.OpenRouterToken, app.docker, app.botRepo, app.log)
}

func (app *App) initKafkaConsumer() {
	app.kafkaConsumer = kafkaconsumer.New(app.botUsecase, &app.cfg.Kafka, app.log)
	app.log.Info("kafka consumer created")
}

func New(envPath string) (*App, error) {
	var app App
	if err := app.initConfig(envPath); err != nil {
		return nil, err
	}
	if err := app.initLogger(); err != nil {
		return nil, err
	}
	if err := app.initDB(); err != nil {
		return nil, err
	}
	if err := app.initDocker(); err != nil {
		return nil, err
	}
	app.initBotRepo()
	app.initBotUsecase()
	app.initKafkaConsumer()

	return &app, nil
}

func (app *App) Run(ctx context.Context) {
	app.log.Info("starting kafka consumer")
	app.kafkaConsumer.Consume(ctx)
}

func (app *App) Stop() {
	app.log.Info("graceful shutdown")

	doneCh := make(chan any, 1)

	go func() {
		defer close(doneCh)

		if err := app.kafkaConsumer.Close(); err != nil {
			app.log.Error("failed to close kafka consumer gracefully", slog.Any("err", err))
		}
		if err := app.botRepo.Close(); err != nil {
			app.log.Error("failed to close db gracefully")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-doneCh:
		app.log.Info("successfully graceful shutdown")
	case <-ctx.Done():
		app.log.Error("failed graceful shutdown, terminate")
	}
}

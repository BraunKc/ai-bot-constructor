package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/braunkc/ai-bot-constructor/orchestrator-service/config"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/application/usecase/bot"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/auth"
	botdomain "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/domain/bot"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/infra/repo/database"
	kafkaproducer "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/infra/repo/kafka/bot/producer"
	botpersistence "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/infra/repo/persistence/bot"
	grpcserver "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/interfaces/grpc"
	"github.com/braunkc/ai-bot-constructor/orchestrator-service/pkg/log"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	grpcServer    *grpc.Server
	db            *gorm.DB
	botRepo       botdomain.BotRepo
	kafkaProducer kafkaproducer.KafkaProducer
	tokenManager  *auth.TokenManager
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
		slog.Error("failed to create config", slog.Any("err", err))
		os.Exit(1)
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

	if app.cfg.Env == config.Develop {
		if err := db.AutoMigrate(&botpersistence.Bot{}); err != nil {
			return fmt.Errorf("failed to migrate tables: %w", err)
		}
	}
	app.db = db
	app.log.Info("database inited")

	return nil
}

func (app *App) initBotRepo() {
	app.botRepo = botpersistence.NewRepo(app.db, app.log)
	app.log.Info("bot repository created")
}

func (app *App) initKafkaProducer() {
	app.kafkaProducer = kafkaproducer.New(&app.cfg.Kafka, app.log)
	app.log.Info("kafka producer created")
}

func (app *App) initBotUsecase() {
	app.botUsecase = botusecase.New(app.botRepo, app.kafkaProducer, app.log)
}

func (app *App) initTokenManager() error {
	tokenManager, err := auth.NewTokenManager(os.Getenv("SECRET_KEY"))
	if err != nil {
		return fmt.Errorf("failed to create token manager: %w", err)
	}
	app.tokenManager = tokenManager
	app.log.Info("token manager created")

	return nil
}

func (app *App) initGRPCServer() {
	app.grpcServer = grpcserver.New(app.botUsecase, app.tokenManager, app.log)
	app.log.Info("grpc server created")
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
	app.initBotRepo()
	app.initKafkaProducer()
	if err := app.initTokenManager(); err != nil {
		return nil, err
	}
	app.initBotUsecase()
	app.initGRPCServer()

	return &app, nil
}

func (app *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", app.cfg.GRPC.Host, app.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to create grpc listener: %w", err)
	}

	app.log.Info("grpc server serve",
		slog.String("host", app.cfg.GRPC.Host),
		slog.String("port", app.cfg.GRPC.Port),
	)
	if err := app.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc server: %w", err)
	}

	return nil
}

func (app *App) Stop() {
	app.log.Info("graceful shutdown")

	doneCh := make(chan any, 1)

	go func() {
		defer close(doneCh)

		app.grpcServer.GracefulStop()
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

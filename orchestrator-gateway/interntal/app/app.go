package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/braunkc/ai-bot-constructor/orchestrator-gateway/config"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/usecase/bot"
	orchestratorgrpc "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/infra/grpc/orchestrator"
	httpserver "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http"
	"github.com/braunkc/ai-bot-constructor/orchestrator-gateway/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	httpServer         *gin.Engine
	botUsecase         botusecase.BotUsecase
	orchestratorClient botusecase.OrchestratorClient
	log                *slog.Logger
	cfg                *config.Config
}

func (app *App) initConfig(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg, err := config.New(os.Getenv("CONFIG_PATH"))
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
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

func (app *App) initOrchestratorClient() error {
	orchestratorClient, err := orchestratorgrpc.NewClient(&app.cfg.GRPC.OrchestratorService, app.log)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator grpc client: %w", err)
	}
	app.orchestratorClient = orchestratorClient
	app.log.Info("orchestrator client created")

	return nil
}

func (app *App) initBotUsecase() {
	app.botUsecase = botusecase.NewBotUsecase(app.orchestratorClient, app.log)
	app.log.Info("bot usecase created")
}

func (app *App) initHTTPServer() {
	app.httpServer = httpserver.New(app.botUsecase)
	app.log.Info("http server created")
}

func New(envPath string) (*App, error) {
	var app App
	if err := app.initConfig(envPath); err != nil {
		return nil, err
	}
	if err := app.initLogger(); err != nil {
		return nil, err
	}
	if err := app.initOrchestratorClient(); err != nil {
		return nil, err
	}
	app.initBotUsecase()
	app.initHTTPServer()

	return &app, nil
}

func (app *App) Run() error {
	app.log.Info("http server run",
		slog.String("host", app.cfg.HTTP.Host),
		slog.String("port", app.cfg.HTTP.Port),
	)
	if err := app.httpServer.Run(fmt.Sprintf("%s:%s", app.cfg.HTTP.Host, app.cfg.HTTP.Port)); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop() {
	app.log.Info("graceful shutdown")

	doneCh := make(chan any, 1)

	go func() {
		defer close(doneCh)

		if err := app.orchestratorClient.Close(); err != nil {
			app.log.Error("failed to close grpc orchestrator conn")
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

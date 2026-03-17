package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/orchestrator-gateway/config"
	botusecase "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/application/usecase/bot"
	orchestratorgrpc "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/infra/grpc/orchestrator"
	httpserver "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/interntal/interfaces/http"
	"github.com/braunkc/ai-bot-constructor/orchestrator-gateway/pkg/log"
)

func main() {
	cfg, err := config.New(".env")
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

	orchesstratorClient, err := orchestratorgrpc.NewClient(&cfg.GRPC.OrchestratorService, log)
	if err != nil {
		log.Error("failed to create orchestrator client")
	}

	botUsecase := botusecase.NewBotUsecase(orchesstratorClient, log)

	server := httpserver.New(botUsecase)

	log.Info("server running")
	if err := server.Run(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
		log.Error("failed to run http server", slog.Any("err", err))
		os.Exit(1)
	}
}

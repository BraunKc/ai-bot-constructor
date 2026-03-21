package main

import (
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/bot-container/config"
	openrouter "github.com/braunkc/ai-bot-constructor/bot-container/internal/infra/open_router"
	tg "github.com/braunkc/ai-bot-constructor/bot-container/internal/infra/telegram"
	"github.com/braunkc/ai-bot-constructor/bot-container/pkg/log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env", slog.Any("err", err))
		os.Exit(1)
	}

	cfg, err := config.New("./config/config.yml")
	if err != nil {
		slog.Error("failed to create config", slog.Any("err", err))
		os.Exit(1)
	}

	logCfg := log.Config{
		Service:    cfg.Logger.Service,
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}
	handler, err := log.NewHandler(&logCfg)
	if err != nil {
		slog.Error("failed to create logger handler", slog.Any("err", err))
		os.Exit(1)
	}
	log := slog.New(handler)
	log.Debug("logger inited")

	openRouter, err := openrouter.New(cfg.OpenRouter.Token, cfg.OpenRouter.Model, log)
	if err != nil {
		log.Error("failed to init open router")
		os.Exit(1)
	}
	log.Info("open router inited")

	tg, err := tg.NewBot(cfg.Telegram.Token, cfg.SystemPrompt, openRouter, log)
	if err != nil {
		log.Error("faliled to init telegram bot", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("telegram bot inited")

	log.Info("bot started")
	tg.Listen()
}

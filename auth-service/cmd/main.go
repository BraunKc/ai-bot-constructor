package main

import (
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/auth-service/config"
	databaserepo "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/database"
	userpersistence "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/persistence/user"
	"github.com/braunkc/ai-bot-constructor/auth-service/pkg/log"
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
	slog.Debug("cfg created")

	logCfg := log.Config{
		Service:    "auth-service",
		OutputType: log.Console,
		Level:      slog.LevelDebug,
	}
	h, err := log.NewHandler(&logCfg)
	if err != nil {
		slog.Error("failed to create log handler", slog.Any("err", err))
		os.Exit(1)
	}
	log := slog.New(h)

	db, err := databaserepo.New(&cfg.DB)
	if err != nil {
		log.Error("failed to create database repository", slog.Any("err", err))
		os.Exit(1)
	}

	if cfg.Env == config.Develop {
		if err := db.AutoMigrate(&userpersistence.User{}); err != nil {
			log.Error("failed to migrate tables", slog.Any("err", err))
			os.Exit(1)
		}
	}
	log.Info("database inited")

	userRepo := userpersistence.NewRepo(db, log)
	log.Info("user repository created")
}

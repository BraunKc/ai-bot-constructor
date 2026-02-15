package main

import (
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/auth-service/config"
	userusecase "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/usecase/user"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/database"
	orchestratorgrpc "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/grpc/orchestrator"
	hasherinfra "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/hasher"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/jwt"
	userpersistence "github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/persistence/user"
	"github.com/braunkc/ai-bot-constructor/auth-service/pkg/log"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	db, err := database.New(&cfg.DB)
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

	orchestratorClient, err := orchestratorgrpc.NewClient(&cfg.GRPC.OrchestratorService, log)
	if err != nil {
		log.Error("failed to create orchestrator grpc client", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("orchestrator client created")

	tokenManager, err := jwt.NewTokenManager(os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Error("failed to create token manager", slog.Any("err", err))
		os.Exit(1)
	}
	log.Info("token manager created")

	userUsecase := userusecase.New(userRepo, orchestratorClient, tokenManager, &hasherinfra.Hasher{Cost: bcrypt.DefaultCost}, log)
	log.Info("user usecase created")
}

package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/braunkc/ai-bot-constructor/auth-service/config"
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

	fmt.Println(cfg)
}

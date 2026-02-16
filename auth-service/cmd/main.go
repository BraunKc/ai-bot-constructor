package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/braunkc/ai-bot-constructor/auth-service/internal/app"
)

func main() {
	app, err := app.New(".env")
	if err != nil {
		slog.Error("failed to init app", slog.Any("err", err))
	}

	go func() {
		if err := app.Run(); err != nil {
			slog.Error("failed to run app", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	app.Stop()
}

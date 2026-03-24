package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/braunkc/ai-bot-constructor/executor-service/internal/app"
)

func main() {
	app, err := app.New(".env")
	if err != nil {
		slog.Error("failed to init app", slog.Any("err", err))
		os.Exit(1)
	}

	kafkaContext, cancel := context.WithCancel(context.Background())
	go app.Run(kafkaContext)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	cancel()
	app.Stop()
}

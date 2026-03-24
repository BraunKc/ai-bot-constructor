package docker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type Docker interface {
	CreateContainer(ctx context.Context, userID, botID uuid.UUID, botName, tgBotToken, openRouterToken, systemPrompt string) error
	StartContainer(ctx context.Context, userID, botID uuid.UUID) error
	StopContainer(ctx context.Context, userID, botID uuid.UUID) error
	RestartContainer(ctx context.Context, userID, botID uuid.UUID) error
	DeleteContainer(ctx context.Context, userID, botID uuid.UUID) error
	Close() error
}

type docker struct {
	apiClient *client.Client
	log       *slog.Logger
}

func New(log *slog.Logger) (Docker, error) {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create apiClient: %w", err)
	}

	return &docker{
		apiClient: apiClient,
		log:       log,
	}, nil
}

func (d *docker) CreateContainer(ctx context.Context, userID, botID uuid.UUID, botName, tgBotToken, openRouterToken, systemPrompt string) error {
	containerName := d.getContainerName(userID, botID)

	d.log.Debug("creating container",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("bot_name", botName),
		slog.String("container_name", containerName),
	)

	_, err := d.apiClient.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err == nil {
		return fmt.Errorf("container already exists: %s", containerName)
	}

	if _, err := d.apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: "ai-bot-container",
		Config: &container.Config{
			Env: []string{
				"CONFIG_PATH=./config/config.yml",
				fmt.Sprintf("TELEGRAM_BOT_TOKEN=%s", tgBotToken),
				fmt.Sprintf("OPEN_ROUTER_TOKEN=%s", openRouterToken),
				fmt.Sprintf("SYSTEM_PROMPT=%s", systemPrompt),
			},
		},
		Name: containerName,
	}); err != nil {
		return err
	}

	return nil
}

func (d *docker) StartContainer(ctx context.Context, userID, botID uuid.UUID) error {
	containerName := d.getContainerName(userID, botID)

	d.log.Debug("starting container",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("container_name", containerName),
	)

	if _, err := d.apiClient.ContainerStart(ctx, containerName, client.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (d *docker) StopContainer(ctx context.Context, userID, botID uuid.UUID) error {
	containerName := d.getContainerName(userID, botID)

	d.log.Debug("stopping container",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("container_name", containerName),
	)

	timeout := 60
	if _, err := d.apiClient.ContainerStop(ctx, containerName, client.ContainerStopOptions{Timeout: &timeout}); err != nil {
		return err
	}

	return nil
}

func (d *docker) RestartContainer(ctx context.Context, userID, botID uuid.UUID) error {
	containerName := d.getContainerName(userID, botID)

	d.log.Debug("restarting container",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("container_name", containerName),
	)

	timeout := 60
	if _, err := d.apiClient.ContainerRestart(ctx, containerName, client.ContainerRestartOptions{Timeout: &timeout}); err != nil {
		return err
	}

	return nil
}

func (d *docker) DeleteContainer(ctx context.Context, userID, botID uuid.UUID) error {
	containerName := d.getContainerName(userID, botID)

	d.log.Debug("deleting container",
		slog.String("user_id", userID.String()),
		slog.String("bot_id", botID.String()),
		slog.String("container_name", containerName),
	)

	if _, err := d.apiClient.ContainerStop(ctx, containerName, client.ContainerStopOptions{}); err != nil {
		d.log.Warn("failed to stop container, deleting force")
	}

	if _, err := d.apiClient.ContainerRemove(ctx, containerName, client.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}

func (d *docker) Close() error {
	return d.apiClient.Close()
}

func (d *docker) getContainerName(userID, botID uuid.UUID) string {
	return fmt.Sprintf("%s-%s", userID.String(), botID.String())
}

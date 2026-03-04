package botdto

import (
	"time"

	"github.com/google/uuid"
)

type Bot struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	BotStatus int32
	Name      string
	LastError string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateBotReq struct {
	Name   string
	ApiKey string
}

type GetBotReq struct {
	ID uuid.UUID
}

type GetAllBotsResp struct {
	Bots []Bot
}

type StopBotReq struct {
	ID uuid.UUID
}

type StopBotsReq struct {
	IDs []uuid.UUID
}

type StopBotsResp struct {
	Bots         []Bot
	AllSucceeded bool
}

type StartBotReq struct {
	ID uuid.UUID
}

type StartBotsReq struct {
	IDs []uuid.UUID
}

type StartBotsResp struct {
	Bots         []Bot
	AllSucceeded bool
}

type RestartBotReq struct {
	ID uuid.UUID
}

type DeleteBotReq struct {
	ID uuid.UUID
}

type DeleteBotsReq struct {
	IDs []uuid.UUID
}

type DeleteBotsResp struct {
	AllSucceeded bool
}

type DeleteAllBotsReq struct {
	UserID uuid.UUID
}

type DeleteAllBotsResp struct {
	AllSucceeded bool
}

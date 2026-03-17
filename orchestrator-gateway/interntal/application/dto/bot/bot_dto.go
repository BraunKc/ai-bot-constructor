package botdto

import (
	"time"

	orchestratorpb "github.com/braunkc/ai-bot-constructor/orchestrator-gateway/api/orchestrator-service/v1"
)

type BotStatus string

const (
	StatusUnknown    BotStatus = "unknown"
	StatusCreating   BotStatus = "creating"
	StatusStarting   BotStatus = "starting"
	StatusRunning    BotStatus = "running"
	StatusStopping   BotStatus = "stopping"
	StatusStopped    BotStatus = "stopped"
	StatusRestarting BotStatus = "restarting"
	StatusDeleting   BotStatus = "deleting"
	StatusError      BotStatus = "error"
)

func ToBotStatus(status orchestratorpb.BotStatus) BotStatus {
	switch status {
	case orchestratorpb.BotStatus_BOT_STATUS_UNKNOWN:
		return StatusUnknown
	case orchestratorpb.BotStatus_BOT_STATUS_CREATING:
		return StatusCreating
	case orchestratorpb.BotStatus_BOT_STATUS_STARTING:
		return StatusStarting
	case orchestratorpb.BotStatus_BOT_STATUS_RUNNING:
		return StatusRunning
	case orchestratorpb.BotStatus_BOT_STATUS_STOPPING:
		return StatusStopping
	case orchestratorpb.BotStatus_BOT_STATUS_STOPPED:
		return StatusStopped
	case orchestratorpb.BotStatus_BOT_STATUS_RESTARTING:
		return StatusRestarting
	case orchestratorpb.BotStatus_BOT_STATUS_DELETING:
		return StatusDeleting
	case orchestratorpb.BotStatus_BOT_STATUS_ERROR:
		return StatusError
	default:
		return StatusUnknown
	}
}

type Bot struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Status    BotStatus `json:"status"`
	LastError string    `json:"last_error"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromProto(pb *orchestratorpb.Bot) *Bot {
	if pb == nil {
		return nil
	}

	dto := &Bot{
		ID:        pb.Id,
		UserID:    pb.UserId,
		Name:      pb.Name,
		Status:    ToBotStatus(pb.Status),
		LastError: pb.LastError,
	}

	if pb.CreatedAt != nil {
		dto.CreatedAt = pb.CreatedAt.AsTime()
	}
	if pb.UpdatedAt != nil {
		dto.UpdatedAt = pb.UpdatedAt.AsTime()
	}

	return dto
}

func FromProtoList(pbs []*orchestratorpb.Bot) []*Bot {
	if pbs == nil {
		return nil
	}

	dtos := make([]*Bot, 0, len(pbs))
	for _, pb := range pbs {
		if dto := FromProto(pb); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

type CreateBotRequest struct {
	Name   string `json:"name"`
	APIKey string `json:"api_key"`
}

type GetBotRequest struct {
	ID string `json:"id"`
}

type IDsRequest struct {
	IDs []string `json:"ids"`
}

type OperationResponse struct {
	AllSucceeded bool   `json:"all_succeeded"`
	Bots         []*Bot `json:"bots"`
}

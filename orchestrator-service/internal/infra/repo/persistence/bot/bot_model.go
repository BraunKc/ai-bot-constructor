package botpersistence

import (
	"time"

	botdomain "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/domain/bot"
	"github.com/google/uuid"
)

type Bot struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	Status       int32     `gorm:"not null;default:1"`
	Name         string    `gorm:"type:varchar(32);not null"`
	SystemPrompt string    `gorm:"type:varchar(1024);not null"`
	ApiKey       string    `gorm:"type:varchar(128);not null; uniqueIndex"`
	LastError    string    `gorm:"type:varchar(256);default:''"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func botDomainToDBModel(bot *botdomain.Bot) *Bot {
	return &Bot{
		ID:           bot.ID(),
		UserID:       bot.UserID(),
		Status:       bot.Status().Int32(),
		Name:         bot.Name().String(),
		SystemPrompt: bot.SystemPrompt(),
		ApiKey:       bot.ApiKey().Raw(),
		LastError:    bot.LastError(),
		CreatedAt:    bot.CreatedAt(),
		UpdatedAt:    bot.UpdatedAt(),
	}
}

func (b *Bot) toDomain() (*botdomain.Bot, error) {
	return botdomain.RestoreBot(b.ID, b.UserID, b.Status, b.Name, b.SystemPrompt, b.ApiKey, b.LastError, b.CreatedAt, b.UpdatedAt)
}

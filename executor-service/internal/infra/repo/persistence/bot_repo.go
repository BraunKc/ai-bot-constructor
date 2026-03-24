package botpersistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	botdomain "github.com/braunkc/ai-bot-constructor/executor-service/internal/domain/bot"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type botRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewBotRepo(db *gorm.DB, log *slog.Logger) botdomain.BotRepo {
	return &botRepo{
		db:  db,
		log: log,
	}
}

func (br *botRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status botdomain.Status) error {
	br.log.Debug("updating bot status",
		slog.String("id", id.String()),
		slog.String("status", status.String()),
	)

	if err := br.db.WithContext(ctx).Model(&Bot{}).Where("id = ?", id).Update("status", status.Int32()).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return botdomain.ErrRecordNotFound
		}

		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (br *botRepo) UpdateError(ctx context.Context, id uuid.UUID, errMsg string) error {
	br.log.Debug("updating bot error",
		slog.String("id", id.String()),
		slog.String("error", errMsg),
	)

	if err := br.db.WithContext(ctx).Model(&Bot{}).Where("id = ?", id).Update("last_error", errMsg).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return botdomain.ErrRecordNotFound
		}

		return fmt.Errorf("failed to update bot status: %w", err)
	}

	return nil
}

func (br *botRepo) Delete(ctx context.Context, id uuid.UUID) error {
	br.log.Debug("deleting bot", slog.String("id", id.String()))

	tx := br.db.WithContext(ctx).Where("id = ?", id).Delete(&Bot{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return botdomain.ErrRecordNotFound
	}

	return nil
}

func (br *botRepo) Close() error {
	sqlDB, err := br.db.DB()
	if err != nil {
		return fmt.Errorf("failed to ger sqlDB: %w", err)
	}

	return sqlDB.Close()
}

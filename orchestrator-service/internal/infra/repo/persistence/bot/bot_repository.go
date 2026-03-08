package botpersistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	botdomain "github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/domain/bot"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type botRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewRepo(db *gorm.DB, log *slog.Logger) botdomain.BotRepo {
	return &botRepo{
		db:  db,
		log: log,
	}
}

func (br *botRepo) Create(ctx context.Context, bot *botdomain.Bot) error {
	br.log.Debug("creating bot",
		slog.String("id", bot.ID().String()),
		slog.String("user_id", bot.UserID().String()),
		slog.String("status", bot.Status().String()),
		slog.String("name", bot.Name().String()),
	)

	if err := br.db.WithContext(ctx).Create(botDomainToDBModel(bot)).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrDuplicatedKey) {
			return botdomain.ErrDuplicatedKey
		}

		return fmt.Errorf("failed to create bot: %w", err)
	}

	return nil
}

func (br *botRepo) GetByID(ctx context.Context, id uuid.UUID) (*botdomain.Bot, error) {
	br.log.Debug("getting bot by id", slog.String("id", id.String()))

	var bot Bot
	if err := br.db.WithContext(ctx).Where("id = ?", id).First(&bot).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return nil, botdomain.ErrRecordNotFound
		}

		return nil, fmt.Errorf("failed to get bot by id: %w", err)
	}

	return bot.toDomain()
}

func (br *botRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*botdomain.Bot, error) {
	br.log.Debug("getting bots by user_id", slog.String("user_id", userID.String()))

	var bots []Bot
	if err := br.db.WithContext(ctx).Where("user_id = ?", userID).Find(&bots).Error; err != nil {
		return nil, fmt.Errorf("failed to get bots by user_id: %w", err)
	}

	domainBots := make([]*botdomain.Bot, 0, len(bots))
	for _, bot := range bots {
		b, err := bot.toDomain()
		if err != nil {
			return nil, err
		}

		domainBots = append(domainBots, b)
	}

	return domainBots, nil
}

func (br *botRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status botdomain.Status) error {
	br.log.Debug("updating bot by id",
		slog.String("id", id.String()),
		slog.String("status", status.String()),
	)

	return br.db.WithContext(ctx).Model(&Bot{}).Where("id = ?", id).Update("status", status.Int32()).Error
}

func (br *botRepo) Close() error {
	sqlDB, err := br.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

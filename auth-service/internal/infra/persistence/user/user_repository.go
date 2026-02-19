package userpersistence

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	userdomain "github.com/braunkc/ai-bot-constructor/auth-service/internal/domain/user"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type userRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewRepo(db *gorm.DB, log *slog.Logger) userdomain.UserRepo {
	return &userRepo{
		db:  db,
		log: log,
	}
}

func (ur *userRepo) Create(ctx context.Context, user *userdomain.User) error {
	ur.log.Debug("creating user",
		slog.String("id", user.ID().String()),
		slog.String("username", user.Username().String()),
	)

	if err := ur.db.WithContext(ctx).Create(userDomainToDBModel(user)).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrDuplicatedKey) {
			return userdomain.ErrDuplicatedKey
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (ur *userRepo) GetByUsername(ctx context.Context, username userdomain.Username) (*userdomain.User, error) {
	ur.log.Debug("getting user by username", slog.String("username", username.String()))

	var user User
	if err := ur.db.WithContext(ctx).Where("username = ?", username.String()).First(&user).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrRecordNotFound
		}

		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user.ToDomain()
}

func (ur *userRepo) Get(ctx context.Context, id uuid.UUID) (*userdomain.User, error) {
	ur.log.Debug("getting user by id", slog.String("id", id.String()))

	var user User
	if err := ur.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrRecordNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user.ToDomain()
}

func (ur *userRepo) UpdateUsername(ctx context.Context, id uuid.UUID, newUsername userdomain.Username) (*userdomain.User, error) {
	ur.log.Debug("updating user username",
		slog.String("id", id.String()),
		slog.String("new_username", newUsername.String()),
	)

	var user User
	if err := ur.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrRecordNotFound) {
			return nil, userdomain.ErrRecordNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user.Username == newUsername.String() {
		return user.ToDomain()
	}

	user.Username = newUsername.String()

	if err := ur.db.WithContext(ctx).Model(&user).Update("username", newUsername.String()).Error; err != nil {
		if errors.Is(postgres.Dialector{}.Translate(err), gorm.ErrDuplicatedKey) {
			return nil, userdomain.ErrDuplicatedKey
		}

		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user.ToDomain()
}

func (ur *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ur.log.Debug("deleting user", slog.String("id", id.String()))

	result := ur.db.WithContext(ctx).Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user by id: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return userdomain.ErrRecordNotFound
	}

	return nil
}

func (ur *userRepo) Close() error {
	sqlDB, err := ur.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

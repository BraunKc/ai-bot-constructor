package userpersistence

import (
	"time"

	userdomain "github.com/braunkc/ai-bot-constructor/auth-service/internal/domain/user"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"type:varchar(32);not null;uniqueIndex"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func userDomainToDBModel(user *userdomain.User) *User {
	return &User{
		ID:           user.ID(),
		Username:     user.Username().String(),
		PasswordHash: user.PasswordHash().String(),
	}
}

func (u *User) ToDomain() (*userdomain.User, error) {
	return userdomain.RestoreUser(u.ID, u.Username, u.PasswordHash)
}

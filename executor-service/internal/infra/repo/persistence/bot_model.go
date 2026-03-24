package botpersistence

import (
	"github.com/google/uuid"
)

type Bot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status    int32     `gorm:"not null"`
	LastError string    `gorm:"type:varchar(256);default:''"`
}

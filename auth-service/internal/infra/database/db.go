package database

import (
	"fmt"
	"time"

	"github.com/braunkc/ai-bot-constructor/auth-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

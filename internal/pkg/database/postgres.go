package database

import (
	"fmt"

	"github.com/OmniscienIT/GolangAPI/config"
	"github.com/OmniscienIT/GolangAPI/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.Category{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderItem{},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

package database

import (
	"fmt"
	"log"

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
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	// Автоматическая миграция (создание таблиц)
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

	log.Println("Database connection established and migrations applied.")
	return db, nil
}

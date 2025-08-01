package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sharer/internal/modules/category"
	"sharer/internal/modules/page"
	"sharer/internal/modules/user"
)

// Config holds database configuration
type Config struct {
	DSN     string
	LogMode logger.LogLevel
}

// NewConnection creates a new GORM database connection
func NewConnection(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogMode),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations for all models
func Migrate(db *gorm.DB) error {
	// Auto-migrate all models
	err := db.AutoMigrate(
		&category.Category{},
		&page.Page{},
		&user.User{}, // Example model, not implemented
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

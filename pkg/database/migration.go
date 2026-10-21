package database

import (
	"github.com/Maulighadge/GolangTemplate/internal/models"
	"github.com/Maulighadge/GolangTemplate/pkg/logger"
	"gorm.io/gorm"
)

// AutoMigrate runs schema auto-migrations for all registered models
func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database auto-migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		logger.Error("Database migration failed: %v", err)
		return err
	}

	logger.Info("Database auto-migration completed successfully")
	return nil
}

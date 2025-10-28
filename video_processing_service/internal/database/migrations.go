package database

import (
	"fmt"
	"log/slog"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

// AutoMigrate runs all database migrations
func AutoMigrate(s Service) error {
	db := s.GetDB()

	slog.Info("🔄 Running database auto-migrations...")

	err := db.AutoMigrate(
		&models.Video{},
		&models.Transcript{},
		&models.Job{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("✅ Database auto-migrations completed successfully.")
	return nil
}

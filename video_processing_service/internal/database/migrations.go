package database

import (
	"fmt"
	"log/slog"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
)

// AutoMigrate runs all database migrations in the correct order
func AutoMigrate(s Service) error {
	db := s.GetDB()

	slog.Info("🔄 Running database auto-migrations...\n")

	slog.Info("📊 Migrating videos table...")
	if err := db.AutoMigrate(&models.Video{}); err != nil {
		return fmt.Errorf("failed to migrate videos table: %w", err)
	}

	slog.Info("📊 Migrating transcripts table...")
	if err := db.AutoMigrate(&models.Transcript{}); err != nil {
		return fmt.Errorf("failed to migrate transcripts table: %w", err)
	}

	slog.Info("✅ Database auto-migrations completed successfully.")
	return nil
}

// DropAllTables drops all tables (use with caution, only for development)
func DropAllTables(s Service) error {
	db := s.GetDB()

	slog.Warn("\n ⚠️  Dropping all database tables... \n")

	if err := db.Migrator().DropTable(&models.Transcript{}); err != nil {
		return fmt.Errorf("failed to drop transcripts table: %w", err)
	}

	if err := db.Migrator().DropTable(&models.Video{}); err != nil {
		return fmt.Errorf("failed to drop videos table: %w", err)
	}

	slog.Info("\n ✅ All tables dropped successfully. \n")
	return nil
}

package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

// RunMigrations executes all necessary database migrations
func (s *service) RunMigrations(ctx context.Context) error {
	log.Println("Running database migrations...")

	if err := s.createVideosTable(ctx); err != nil {
		return fmt.Errorf("failed to create videos table: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createVideosTable creates the videos table if it doesn't exist
func (s *service) createVideosTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS videos (
			id SERIAL PRIMARY KEY,
			video_id VARCHAR(255) NOT NULL,
			platform VARCHAR(50) NOT NULL,
			original_url TEXT NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(video_id, platform)
		);

		CREATE INDEX IF NOT EXISTS idx_videos_video_id ON videos(video_id);
		CREATE INDEX IF NOT EXISTS idx_videos_status ON videos(status);
		CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos(created_at DESC);
	`

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	log.Println("✓ Videos table created/verified")
	return nil
}

// AutoMigrate is a convenience function to run migrations on the service
func AutoMigrate(s Service) error {
	ctx := context.Background()

	// Type assert to access the concrete service type
	if svc, ok := s.(*service); ok {
		return svc.RunMigrations(ctx)
	}

	return fmt.Errorf("unable to run migrations: invalid service type")
}

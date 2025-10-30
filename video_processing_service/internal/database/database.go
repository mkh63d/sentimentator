package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Oskarowski/sentimentator/video_processing_service/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *gorm.DB

	// Video operations
	SaveVideo(ctx context.Context, video *models.Video) error
	GetVideoByID(ctx context.Context, id uint) (*models.Video, error)
	GetVideoByVideoID(ctx context.Context, videoID, platform string) (*models.Video, error)
	UpdateVideoStatus(ctx context.Context, id uint, status string) error

	// Transcript operations
	CreateTranscript(ctx context.Context, transcript *models.Transcript) error
	GetTranscriptByID(ctx context.Context, id uint) (*models.Transcript, error)
	GetTranscriptByVideoID(ctx context.Context, videoID uint, language string) (*models.Transcript, error)
	HasTranscript(ctx context.Context, videoID uint, language string) (bool, *models.Transcript, error)
	UpdateTranscriptStatus(ctx context.Context, id uint, status models.TranscriptStatus) error

	// Video list operations
	ListVideos(ctx context.Context, limit, offset int) ([]models.Video, int64, error)
}

type service struct {
	db *gorm.DB
}

var (
	database   = os.Getenv("SENTIMENTATOR_VPS_DB_DATABASE")
	password   = os.Getenv("SENTIMENTATOR_VPS_DB_PASSWORD")
	username   = os.Getenv("SENTIMENTATOR_VPS_DB_USERNAME")
	port       = os.Getenv("SENTIMENTATOR_VPS_DB_PORT")
	host       = os.Getenv("SENTIMENTATOR_VPS_DB_HOST")
	schema     = os.Getenv("SENTIMENTATOR_VPS_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s TimeZone=UTC",
		host, username, password, database, port, schema,
	)

	gormLogLevelStr := os.Getenv("SENTIMENTATOR_GORM_LOG_LEVEL")
	if gormLogLevelStr == "" {
		gormLogLevelStr = os.Getenv("GORM_LOG_LEVEL")
	}

	var gormLogLevel logger.LogLevel = logger.Info
	switch strings.ToLower(gormLogLevelStr) {
	case "silent":
		gormLogLevel = logger.Silent
	case "error":
		gormLogLevel = logger.Error
	case "warn", "warning":
		gormLogLevel = logger.Warn
	case "info":
		gormLogLevel = logger.Info
	case "":
		gormLogLevel = logger.Error
	default:
		gormLogLevel = logger.Error
		log.Printf("unknown SENTIMENTATOR_GORM_LOG_LEVEL='%s', defaulting to 'error'", gormLogLevelStr)
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	dbInstance = &service{db: db}
	return dbInstance
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db error: %v", err)
		return stats
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	slog.Info("Disconnected from database", "database", database)
	return sqlDB.Close()
}

func (s *service) SaveVideo(ctx context.Context, video *models.Video) error {
	return s.db.WithContext(ctx).Create(video).Error
}

func (s *service) GetVideoByID(ctx context.Context, id uint) (*models.Video, error) {
	var video models.Video
	err := s.db.WithContext(ctx).
		Preload("Transcripts").
		First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (s *service) GetVideoByVideoID(ctx context.Context, videoID, platform string) (*models.Video, error) {
	var video models.Video
	err := s.db.WithContext(ctx).
		Where("video_id = ? AND platform = ?", videoID, platform).
		Preload("Transcripts").
		First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (s *service) UpdateVideoStatus(ctx context.Context, id uint, status string) error {
	return s.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Transcript operations

func (s *service) CreateTranscript(ctx context.Context, transcript *models.Transcript) error {
	return s.db.WithContext(ctx).Create(transcript).Error
}

func (s *service) GetTranscriptByID(ctx context.Context, id uint) (*models.Transcript, error) {
	var transcript models.Transcript
	err := s.db.WithContext(ctx).
		First(&transcript, id).Error
	if err != nil {
		return nil, err
	}
	return &transcript, nil
}

func (s *service) GetTranscriptByVideoID(ctx context.Context, videoID uint, language string) (*models.Transcript, error) {
	var transcript models.Transcript
	err := s.db.WithContext(ctx).
		Where("video_id = ? AND language = ? AND status IN ?",
			videoID,
			language,
			[]models.TranscriptStatus{
				models.StatusCompleted,
				models.StatusProcessing,
				models.StatusPending,
			}).
		Order("created_at DESC").
		First(&transcript).Error

	if err != nil {
		return nil, err
	}

	return &transcript, nil
}

// HasTranscript checks if a video already has a transcript for the given language
func (s *service) HasTranscript(ctx context.Context, videoID uint, language string) (bool, *models.Transcript, error) {
	transcript, err := s.GetTranscriptByVideoID(ctx, videoID, language)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, transcript, nil
}

func (s *service) UpdateTranscriptStatus(ctx context.Context, id uint, status models.TranscriptStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.StatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	return s.db.WithContext(ctx).
		Model(&models.Transcript{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (s *service) ListVideos(ctx context.Context, limit, offset int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	if err := s.db.WithContext(ctx).Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := s.db.WithContext(ctx).
		Preload("Transcripts").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error

	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

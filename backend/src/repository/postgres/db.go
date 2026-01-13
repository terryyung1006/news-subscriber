package postgres

import (
	"fmt"
	"news-subscriber-core/src/config"
	"news-subscriber-core/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	return db.AutoMigrate(
		&models.User{},
		&models.InviteCode{},
	)
}

// Reset drops all tables and recreates them (use with caution!)
func (db *DB) Reset() error {
	if err := db.Migrator().DropTable(&models.InviteCode{}, &models.User{}); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	return db.Migrate()
}

// Seed initial invite codes
func (db *DB) Seed() error {
	codes := []string{"WELCOME2024", "BETA_USER_001", "EARLY_ACCESS"}

	for _, code := range codes {
		var count int64
		db.Model(&models.InviteCode{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			if err := db.Create(&models.InviteCode{Code: code}).Error; err != nil {
				return fmt.Errorf("failed to seed invite code %s: %w", code, err)
			}
		}
	}
	return nil
}

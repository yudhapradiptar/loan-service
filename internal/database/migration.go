package database

import (
	"fmt"
	"log"

	"loan-service/internal/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RunMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func RollbackMigration(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Database migration rollback completed successfully")
	return nil
}

package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDatabase реализует интерфейс Database для PostgreSQL
type PostgresDatabase struct {
	db *gorm.DB
}

// NewPostgresDatabase создает новое подключение к PostgreSQL базе данных
func NewPostgresDatabase(config DatabaseConfig) (Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
		config.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully created new PostgreSQL database connection!")
	
	return &PostgresDatabase{
		db: db,
	}, nil
}

// GetDB возвращает экземпляр *gorm.DB для использования в репозиториях
func (p *PostgresDatabase) GetDB() *gorm.DB {
	return p.db
}

// Close закрывает соединение с базой данных
func (p *PostgresDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	
	log.Println("Database connection closed successfully")
	return nil
}

// Ping проверяет соединение с базой данных
func (p *PostgresDatabase) Ping() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	return nil
}

package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresConfig содержит конфигурацию для подключения к PostgreSQL
type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

// postgresDatabase реализует интерфейс Database для PostgreSQL
type postgresDatabase struct {
	db *gorm.DB
}

// NewPostgresDatabase создает новое подключение к PostgreSQL
func NewPostgresDatabase(config PostgresConfig) (Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode, config.TimeZone)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("[DB] Successfully connected to PostgreSQL")

	return &postgresDatabase{
		db: gormDB,
	}, nil
}

// GetDB возвращает экземпляр *gorm.DB для использования в репозиториях
func (p *postgresDatabase) GetDB() *gorm.DB {
	return p.db
}

// Close закрывает соединение с базой данных
func (p *postgresDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Ping проверяет соединение с базой данных
func (p *postgresDatabase) Ping() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Ping()
}

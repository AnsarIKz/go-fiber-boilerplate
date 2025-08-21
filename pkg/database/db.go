package database

import (
	"nodabackend/pkg/env"
	"sync"

	"gorm.io/gorm"
)

var (
	dbOnce sync.Once
	db     Database
	errDB  error
)

// Database интерфейс для работы с базой данных
type Database interface {
	// GetDB возвращает экземпляр *gorm.DB для использования в репозиториях
	GetDB() *gorm.DB
	
	// Close закрывает соединение с базой данных
	Close() error
	
	// Ping проверяет соединение с базой данных
	Ping() error
}

// DatabaseConfig содержит конфигурацию для подключения к базе данных
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}



// NewDatabase создает новое подключение к базе данных
// Конфигурация считывается из переменных окружения или используются значения по умолчанию
func NewDatabase() (Database, error) {
	dbOnce.Do(func() {
		config := DatabaseConfig{
			Host:     env.GetEnv("DB_HOST", "localhost"),
			User:     env.GetEnv("DB_USER", "postgres"),
			Password: env.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   env.GetEnv("DB_NAME", "noda"),
			Port:     env.GetEnv("DB_PORT", "5432"),
			SSLMode:  env.GetEnv("DB_SSLMODE", "disable"),
			TimeZone: env.GetEnv("DB_TIMEZONE", "Asia/Astana"),
		}

		// По умолчанию используем PostgreSQL
		db, errDB = NewPostgresDatabase(config)
	})

	return db, errDB
}

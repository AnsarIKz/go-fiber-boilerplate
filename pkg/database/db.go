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
	GetDB() *gorm.DB
	Close() error
	Ping() error
}




// NewDatabase создает новое подключение к базе данных
// Конфигурация считывается из переменных окружения или используются значения по умолчанию
func NewDatabase() (Database, error) {
	dbOnce.Do(func() {
		config := PostgresConfig{
			Host:     env.GetEnvOrDefault("DB_HOST", "localhost"),
			User:     env.GetEnvOrDefault("DB_USER", "postgres"),
			Password: env.GetEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   env.GetEnvOrDefault("DB_NAME", "noda"),
			Port:     env.GetEnvOrDefault("DB_PORT", "5432"),
			SSLMode:  env.GetEnvOrDefault("DB_SSLMODE", "disable"),
			TimeZone: env.GetEnvOrDefault("DB_TIMEZONE", "Asia/Astana"),
		}

		db, errDB = NewPostgresDatabase(config)
	})

	return db, errDB
}

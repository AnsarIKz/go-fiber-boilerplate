package otp

import (
	"context"
	"fmt"
	"log"
	"nodabackend/pkg/database"
	"nodabackend/pkg/redis"
	"time"
)

func Example() {
	// --- Инициализация ---
	// В реальном приложении это будет делаться один раз в main.go

	// 1. Подключаемся к PostgreSQL
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// 2. Подключаемся к Redis
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	// 3. Создаем OTP сервис, передавая ему Redis клиент
	service := NewOTPService(redisClient, nil)
	// --- Конец инициализации ---

	ctx := context.Background()

	// Пример 1: Генерация OTP для входа
	userID := "user123"
	otpType := OTPTypeLogin

	// Генерируем OTP
	generatedOTP, err := service.GenerateOTP(ctx, userID, otpType)
	if err != nil {
		log.Fatalf("Failed to generate OTP: %v", err)
	}

	fmt.Printf("Generated OTP for user %s: %s\n", userID, generatedOTP)

	// Имитируем отправку OTP пользователю
	fmt.Printf("Sending OTP %s to user %s via SMS/Email...\n", generatedOTP, userID)

	// Пример 2: Проверка валидного OTP
	userOTP := generatedOTP // Пользователь ввел правильный OTP
	isValid, err := service.ValidateOTP(ctx, userID, otpType, userOTP)
	if err != nil {
		log.Fatalf("Failed to validate OTP: %v", err)
	}

	if isValid {
		fmt.Println("✅ OTP is valid! User authenticated successfully.")
	} else {
		fmt.Println("❌ OTP is invalid!")
	}
}

// ExampleWithCustomConfig показывает использование с кастомной конфигурацией
func ExampleWithCustomConfig() {
	// --- Инициализация ---
	// 1. Подключаемся к Redis
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()
	// --- Конец инициализации ---

	// Создаем кастомную конфигурацию для OTP
	config := &OTPConfig{
		OTPLength:   8,
		OTPExpiry:   5 * time.Minute,
		MaxAttempts: 5,
	}

	// Создаем сервис с кастомной конфигурацией
	service := NewOTPService(redisClient, config)

	ctx := context.Background()

	// Используем сервис как обычно
	otp, err := service.GenerateOTP(ctx, "user456", OTPTypeRegister)
	if err != nil {
		log.Fatalf("Failed to generate OTP: %v", err)
	}

	fmt.Printf("Generated OTP with custom config: %s\n", otp)
}







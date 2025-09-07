package otp

import (
	"time"
)

// OTPConfig конфигурация OTP сервиса
type OTPConfig struct {
	OTPLength   int
	OTPExpiry   time.Duration
	MaxAttempts int
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *OTPConfig {
	return &OTPConfig{
		OTPLength:   6,
		OTPExpiry:   10 * time.Minute,
		MaxAttempts: 3,
	}
}







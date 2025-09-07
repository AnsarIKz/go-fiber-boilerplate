package otp

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// OTPTypeLogin для входа в систему
	OTPTypeLogin = "login"
	// OTPTypeRegister для регистрации
	OTPTypeRegister = "register"
	// OTPTypeResetPassword для сброса пароля
	OTPTypeResetPassword = "reset_password"
	// OTPTypeVerifyEmail для верификации email
	OTPTypeVerifyEmail = "verify_email"
	// OTPTypeVerifyPhone для верификации телефона
	OTPTypeVerifyPhone = "verify_phone"
)

// OTPData структура для хранения информации об OTP
type OTPData struct {
	HashedOTP string    `json:"hashed_otp"`
	Salt      string    `json:"salt"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

// OTPService интерфейс для работы с OTP
type OTPService interface {
	GenerateOTP(ctx context.Context, userID, otpType string) (string, error)
	ValidateOTP(ctx context.Context, userID, otpType, otp string) (bool, error)
	GetOTPInfo(ctx context.Context, userID, otpType string) (*OTPData, error)
	DeleteOTP(ctx context.Context, userID, otpType string) error
}

// DefaultOTPService реализация OTP сервиса
type DefaultOTPService struct {
	redisClient *redis.Client
	config      *OTPConfig
}

// NewOTPService создает новый OTP сервис
func NewOTPService(redisClient *redis.Client, config *OTPConfig) OTPService {
	if config == nil {
		config = DefaultConfig()
	}
	return &DefaultOTPService{
		redisClient: redisClient,
		config:      config,
	}
}

// generateRandomOTP генерирует случайный OTP из английских букв и цифр
func (s *DefaultOTPService) generateRandomOTP() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := s.config.OTPLength

	otp := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}
		otp[i] = charset[randomIndex.Int64()]
	}

	return string(otp), nil
}

// makeKey создает ключ для Redis
func (s *DefaultOTPService) makeKey(userID, otpType string) string {
	return fmt.Sprintf("otp:%s:%s", userID, otpType)
}

// GenerateOTP генерирует новый OTP для пользователя
func (s *DefaultOTPService) GenerateOTP(ctx context.Context, userID, otpType string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}

	if otpType == "" {
		return "", errors.New("OTP type is required")
	}

	// Генерируем OTP
	otp, err := s.generateRandomOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Генерируем соль
	salt, err := generateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Хэшируем OTP
	hashedOTP := hashOTP(otp, salt)

	// Создаем данные для хранения
	otpData := &OTPData{
		HashedOTP: hashedOTP,
		Salt:      salt,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	// Сохраняем в Redis
	key := s.makeKey(userID, otpType)
	data := fmt.Sprintf("%s:%s:%d:%d",
		otpData.HashedOTP,
		otpData.Salt,
		otpData.Attempts,
		otpData.CreatedAt.Unix(),
	)

	err = s.redisClient.Set(ctx, key, data, s.config.OTPExpiry).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store OTP in Redis: %w", err)
	}

	return otp, nil
}

// ValidateOTP проверяет OTP
func (s *DefaultOTPService) ValidateOTP(ctx context.Context, userID, otpType, otp string) (bool, error) {
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	if otpType == "" {
		return false, errors.New("OTP type is required")
	}

	if otp == "" {
		return false, errors.New("OTP is required")
	}

	if len(otp) != s.config.OTPLength {
		return false, errors.New("invalid OTP length")
	}

	key := s.makeKey(userID, otpType)

	// Получаем данные из Redis
	data, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, errors.New("OTP not found or expired")
	}
	if err != nil {
		return false, fmt.Errorf("failed to get OTP from Redis: %w", err)
	}

	// Разбираем данные
	var storedHashedOTP, salt string
	var attempts int
	var createdAt int64

	_, err = fmt.Sscanf(data, "%s:%s:%d:%d", &storedHashedOTP, &salt, &attempts, &createdAt)
	if err != nil {
		return false, fmt.Errorf("failed to parse stored OTP data: %w", err)
	}

	// Проверяем количество попыток
	if attempts >= s.config.MaxAttempts {
		// Удаляем OTP после превышения лимита попыток
		s.redisClient.Del(ctx, key)
		return false, errors.New("maximum attempts exceeded")
	}

	// Хэшируем введенный OTP с той же солью
	inputHashedOTP := hashOTP(otp, salt)

	// Проверяем совпадение
	if inputHashedOTP != storedHashedOTP {
		// Увеличиваем счетчик попыток
		attempts++
		newData := fmt.Sprintf("%s:%s:%d:%d", storedHashedOTP, salt, attempts, createdAt)

		// Обновляем данные в Redis
		ttl := s.redisClient.TTL(ctx, key).Val()
		s.redisClient.Set(ctx, key, newData, ttl)

		return false, errors.New("invalid OTP")
	}

	// OTP верный, удаляем его из Redis
	s.redisClient.Del(ctx, key)

	return true, nil
}

// GetOTPInfo получает информацию об OTP (для отладки)
func (s *DefaultOTPService) GetOTPInfo(ctx context.Context, userID, otpType string) (*OTPData, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if otpType == "" {
		return nil, errors.New("OTP type is required")
	}

	key := s.makeKey(userID, otpType)

	data, err := s.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("OTP not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get OTP from Redis: %w", err)
	}

	var hashedOTP, salt string
	var attempts int
	var createdAt int64

	_, err = fmt.Sscanf(data, "%s:%s:%d:%d", &hashedOTP, &salt, &attempts, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stored OTP data: %w", err)
	}

	return &OTPData{
		HashedOTP: hashedOTP,
		Salt:      salt,
		Attempts:  attempts,
		CreatedAt: time.Unix(createdAt, 0),
	}, nil
}

// DeleteOTP удаляет OTP из Redis
func (s *DefaultOTPService) DeleteOTP(ctx context.Context, userID, otpType string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	if otpType == "" {
		return errors.New("OTP type is required")
	}

	key := s.makeKey(userID, otpType)
	return s.redisClient.Del(ctx, key).Err()
}

// generateSalt генерирует случайную соль
func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// hashOTP хэширует OTP с солью для безопасного хранения
func hashOTP(otp, salt string) string {
	data := otp + salt
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}







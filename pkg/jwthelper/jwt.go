package jwthelper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims представляет JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTHelper содержит методы для работы с JWT токенами
type JWTHelper struct {
	secretKey []byte
}

// NewJWTHelper создает новый JWT helper
func NewJWTHelper() *JWTHelper {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // для разработки, в продакшене обязательно через переменную окружения
	}
	return &JWTHelper{
		secretKey: []byte(secret),
	}
}

// GenerateToken создает новый JWT токен
func (j *JWTHelper) GenerateToken(userID uint, phone string, role string) (string, error) {
	// Устанавливаем время жизни токена - 24 часа
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken проверяет и парсит JWT токен
func (j *JWTHelper) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken создает новый токен на основе старого
func (j *JWTHelper) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Создаем новый токен с теми же данными пользователя
	return j.GenerateToken(claims.UserID, claims.Phone, claims.Role)
}


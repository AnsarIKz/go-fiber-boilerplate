package middleware

import (
	"nodabackend/pkg/jwthelper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware middleware для проверки JWT токенов
type AuthMiddleware struct {
	jwtHelper *jwthelper.JWTHelper
}

// NewAuthMiddleware создает новый auth middleware
func NewAuthMiddleware(jwtHelper *jwthelper.JWTHelper) *AuthMiddleware {
	return &AuthMiddleware{
		jwtHelper: jwtHelper,
	}
}

// RequireAuth middleware для защищенных маршрутов
func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	// Получаем токен из заголовка Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Проверяем формат Bearer токена
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	tokenString := tokenParts[1]

	// Валидируем токен
	claims, err := m.jwtHelper.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Сохраняем данные пользователя в контексте для использования в handlers
	c.Locals("userID", claims.UserID)
	c.Locals("phone", claims.Phone)

	return c.Next()
}

// OptionalAuth middleware для маршрутов где авторизация опциональна
func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			claims, err := m.jwtHelper.ValidateToken(tokenParts[1])
			if err == nil {
				c.Locals("userID", claims.UserID)
				c.Locals("phone", claims.Phone)
			}
		}
	}
	return c.Next()
}

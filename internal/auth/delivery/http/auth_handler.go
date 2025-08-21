package http

import (
	"nodabackend/internal/auth/usecase"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler HTTP обработчики для аутентификации
type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler создает новый handler
func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: uc}
}

// LoginRequest структура для запроса входа
type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// RegisterRequest структура для запроса регистрации
type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Login обрабатывает вход пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.authUseCase.AuthenticateUser(req.Phone, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

// Register обрабатывает регистрацию пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.authUseCase.RegisterUser(req.Phone, req.Password, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

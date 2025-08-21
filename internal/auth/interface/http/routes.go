package http

import (
	"nodabackend/internal/auth/middleware"
	"nodabackend/internal/auth/repository"
	"nodabackend/internal/auth/usecase"
	"nodabackend/pkg/jwthelper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes настраивает маршруты аутентификации и возвращает auth middleware
func RegisterRoutes(api fiber.Router, db *gorm.DB, jwtHelper *jwthelper.JWTHelper) *middleware.AuthMiddleware {
	authRepo := repository.NewUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(authRepo, jwtHelper)
	authHandler := NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(jwtHelper)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	auth.Get("/me", authMiddleware.RequireAuth, authHandler.Me)

	return authMiddleware
}

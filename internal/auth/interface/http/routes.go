package http

import (
	"nodabackend/internal/auth/middleware"
	"nodabackend/internal/auth/repository"
	"nodabackend/internal/auth/usecase"
	"nodabackend/pkg/jwthelper"
	"nodabackend/pkg/mailer"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes настраивает маршруты аутентификации и возвращает auth middleware
func RegisterRoutes(api fiber.Router, db *gorm.DB, jwtHelper *jwthelper.JWTHelper) *middleware.AuthMiddleware {
	authRepo := repository.NewUserRepository(db)
	mailer := mailer.NewSMTPMailerFromEnv()
	authUseCase := usecase.NewAuthUseCase(authRepo, jwtHelper, mailer)
	authHandler := NewAuthHandler(authUseCase)
	authMiddleware := middleware.NewAuthMiddleware(jwtHelper)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	auth.Get("/me", authMiddleware.RequireAuth, authHandler.Me)

	// Пример защищенного маршрута для админа
	// Сначала проверяем токен (RequireAuth), потом роль (RequireRole)
	admin := api.Group("/admin", authMiddleware.RequireAuth, authMiddleware.RequireRole("admin"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the admin dashboard!",
		})
	})

	return authMiddleware
}

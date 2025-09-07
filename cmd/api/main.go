package main

import (
	"log"
	"nodabackend/internal/auth/interface/http"
	"nodabackend/pkg/database"
	"nodabackend/pkg/jwthelper"

	// "nodabackend/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Подключение к базе данных (PostgreSQL)
	dbInstance, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db := dbInstance.GetDB()

	// 2. Подключение к Redis
	// redisClient, err := redis.NewClient()
	// if err != nil {
	// 	log.Fatal("Failed to connect to Redis:", err)
	// }

	// 3. Инициализация JWT helper
	jwtHelper := jwthelper.NewJWTHelper()

	// 4. Настройка HTTP сервера
	app := fiber.New()

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Noda Backend API",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Auth routes
	// TODO: Передать redisClient в RegisterRoutes, когда он понадобится
	authMiddleware := http.RegisterRoutes(api, db, jwtHelper)

	// Example protected route
	protected := api.Group("/protected", authMiddleware.RequireAuth)
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		phone := c.Locals("phone").(string)
		
		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"user_id": userID,
			"phone":   phone,
		})
	})

	log.Println("Server starting on :3000")
	app.Listen(":3000")
}
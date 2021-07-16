package middleware

import (
	"github.com/SemmiDev/fiber-go-blog/infrastructure/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := auth.TokenValid(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Next()
	}
}

func CORSMiddleware(app *fiber.App) {
	app.Use(cors.New())
}

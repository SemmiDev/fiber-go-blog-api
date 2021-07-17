package middleware

import (
	auth2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := auth2.TokenValid(c)
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

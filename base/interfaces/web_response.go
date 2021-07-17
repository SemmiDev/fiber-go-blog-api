package interfaces

import "github.com/gofiber/fiber/v2"

func ERROR(c *fiber.Ctx, status int, success bool, message interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": success,
		"message": message,
	})
}

func SUCCESS(c *fiber.Ctx, status int, success bool, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": success,
		"data":    data,
	})
}

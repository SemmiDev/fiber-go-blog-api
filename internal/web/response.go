package web

import (
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/gofiber/fiber/v2"
)

func MarshalPayload(c *fiber.Ctx, code int, payload interface{}) error {
	c.Set("Content-Type", "application/json")
	c.Status(code)
	return c.JSON(payload)
}

func MarshalError(c *fiber.Ctx, code int, err error) error {
	c.Set("Content-Type", "application/json")
	c.Status(code)
	return c.JSON(model.ErrorResponse{Message: err.Error()})
}
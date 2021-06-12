package handler

import (
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/service"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/validation"
	"github.com/SemmiDev/fiber-go-blog/internal/web"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService}
}

type authHandler struct {
	authService service.AuthService
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	var req model.AuthRequest
	err := c.BodyParser(&req)
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	res, err := h.authService.Login(c.Context(), req)
	if err != nil {
		switch err {
		case constant.ErrEmailNotRegistered, constant.ErrWrongPassword:
			return web.MarshalError(c, fiber.StatusUnauthorized, err)
		default:
			return web.MarshalError(c, fiber.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, fiber.StatusOK, res)
}
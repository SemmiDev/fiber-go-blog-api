package handler

import (
	"errors"
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/service"
	"github.com/SemmiDev/fiber-go-blog/internal/auth"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/validation"
	"github.com/SemmiDev/fiber-go-blog/internal/web"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type AccountHandler interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{accountService}
}

type accountHandler struct {
	accountService service.AccountService
}

func (h *accountHandler) Create(c *fiber.Ctx) error {
	var req model.AccountCreateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, err)
	}

	res, err := h.accountService.Create(c, req)
	if err != nil {
		switch err {
		case constant.ErrEmailRegistered:
			return web.MarshalError(c, http.StatusConflict, err)
		default:
			return web.MarshalError(c, http.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, http.StatusCreated, res)
}

func (h *accountHandler) List(c *fiber.Ctx) error {
	limit, offset, err := web.GetPagination(c)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, err)
	}

	req := model.AccountListRequest{
		Limit:  limit,
		Offset: offset,
		Name:   c.Query("name"),
	}

	res, err := h.accountService.List(c, req)
	if err != nil {
		return web.MarshalError(c, http.StatusInternalServerError, err)
	}

	return web.MarshalPayload(c, http.StatusOK, res)
}

func (h *accountHandler) Get(c *fiber.Ctx) error {
	id, err := web.GetUrlPathInt64(c.Params("account_id"))
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, err)
	}

	req := model.AccountGetRequest{ID: id}
	res, err := h.accountService.Get(c, req)
	if err != nil {
		switch err {
		case constant.ErrAccountNotFound:
			return web.MarshalError(c, http.StatusNotFound, err)
		default:
			return web.MarshalError(c, http.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, http.StatusOK, res)
}

func (h *accountHandler) Update(c *fiber.Ctx) error {
	// authorize request
	id, err := authorizeAccount(c)
	if err != nil {
		return err
	}

	log.Println(id)

	req := model.AccountUpdateRequest{ID: id}
	err = c.BodyParser(&req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, err)
	}

	res, err := h.accountService.Update(c, req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			return web.MarshalError(c, http.StatusUnauthorized, err)
		case constant.ErrEmailRegistered:
			return web.MarshalError(c, http.StatusConflict, err)
		case constant.ErrAccountNotFound:
			return web.MarshalError(c, http.StatusNotFound, err)
		default:
			return web.MarshalError(c, http.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, http.StatusOK, res)
}

func (h *accountHandler) UpdatePassword(c *fiber.Ctx) error {
	// authorize request
	id, err := authorizeAccount(c)
	if err != nil {
		return err
	}

	req := model.AccountPasswordUpdateRequest{ID: id}
	err = c.BodyParser(&req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		return web.MarshalError(c, http.StatusBadRequest, err)
	}

	res, err := h.accountService.UpdatePassword(c, req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized, constant.ErrWrongPassword:
			return web.MarshalError(c, http.StatusUnauthorized, err)
		case constant.ErrAccountNotFound:
			return web.MarshalError(c, http.StatusNotFound, err)
		default:
			return web.MarshalError(c, http.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, http.StatusOK, res)
}

func (h *accountHandler) Delete(c *fiber.Ctx) error {
	// authorize request
	id, err := authorizeAccount(c)
	if err != nil {
		return err
	}

	req := model.AccountDeleteRequest{ID: id}
	err = h.accountService.Delete(c, req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			return web.MarshalError(c, http.StatusUnauthorized, err)
		case constant.ErrAccountNotFound:
			return web.MarshalError(c, http.StatusNotFound, err)
		default:
			return web.MarshalError(c, http.StatusInternalServerError, err)
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}


func authorizeAccount(c *fiber.Ctx) (int64, error) {
	//Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c)
	if err != nil {
		return 0, web.MarshalError(c, fiber.StatusUnauthorized, err)
	}

	id, err := web.GetUrlPathInt64(c.Params("account_id"))
	if err != nil {
		return 0, web.MarshalError(c, http.StatusBadRequest, err)
	}

	// check if user_id param == user_id in extracted in jwt
	if uid != id {
		return 0, web.MarshalError(c, fiber.StatusUnauthorized, errors.New(http.StatusText(fiber.StatusUnauthorized)))
	}

	return id, nil
}
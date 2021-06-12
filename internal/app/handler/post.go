package handler

import (
	"github.com/SemmiDev/fiber-go-blog/internal/app/model"
	"github.com/SemmiDev/fiber-go-blog/internal/app/service"
	"github.com/SemmiDev/fiber-go-blog/internal/constant"
	"github.com/SemmiDev/fiber-go-blog/internal/validation"
	"github.com/SemmiDev/fiber-go-blog/internal/web"
	"github.com/gofiber/fiber/v2"
)

type PostHandler interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewPostHandler(postService service.PostService) PostHandler {
	return &postHandler{postService}
}

type postHandler struct {
	postService service.PostService
}

func (h *postHandler) Create(c *fiber.Ctx) error {
	var req model.PostCreateRequest
	err := c.BodyParser(&req)
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	res, err := h.postService.Create(c, req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			return web.MarshalError(c, fiber.StatusUnauthorized, err)
		default:
			return web.MarshalError(c, fiber.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, fiber.StatusCreated, res)
}

func (h *postHandler) List(c *fiber.Ctx) error {
	limit, offset, err := web.GetPagination(c)
	if err != nil {
		web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	req := model.PostListRequest{
		Limit:  limit,
		Offset: offset,
		Title:  c.Query("title"),
	}

	res, err := h.postService.List(c, req)
	if err != nil {
		web.MarshalError(c, fiber.StatusInternalServerError, err)
	}

	return web.MarshalPayload(c, fiber.StatusOK, res)
}

func (h *postHandler) Get(c *fiber.Ctx) error {
	id, err := web.GetUrlPathInt64(c.Params("post_id"))
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	req := model.PostGetRequest{ID: id}
	res, err := h.postService.Get(c, req)
	if err != nil {
		switch err {
		case constant.ErrPostNotFound:
			return web.MarshalError(c, fiber.StatusNotFound, err)
		default:
			return web.MarshalError(c, fiber.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, fiber.StatusOK, res)
}

func (h *postHandler) Update(c *fiber.Ctx) error {
	id, err := web.GetUrlPathInt64(c.Params("post_id"))
	if err != nil {
		return web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	req := model.PostUpdateRequest{ID: id}
	err = c.BodyParser(&req)
	if err != nil {
		web.MarshalError(c, fiber.StatusBadRequest, constant.ErrRequestBody)
	}

	err = validation.Struct(req)
	if err != nil {
		web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	res, err := h.postService.Update(c, req)
	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			return web.MarshalError(c, fiber.StatusUnauthorized, err)
		case constant.ErrPostNotFound:
			return web.MarshalError(c, fiber.StatusNotFound, err)
		default:
			return web.MarshalError(c, fiber.StatusInternalServerError, err)
		}
	}

	return web.MarshalPayload(c, fiber.StatusOK, res)
}

func (h *postHandler) Delete(c *fiber.Ctx) error {
	id, err := web.GetUrlPathInt64(c.Params("post_id"))
	if err != nil {
		web.MarshalError(c, fiber.StatusBadRequest, err)
	}

	req := model.PostDeleteRequest{ID: id}
	err = h.postService.Delete(c, req)

	if err != nil {
		switch err {
		case constant.ErrUnauthorized:
			return web.MarshalError(c, fiber.StatusUnauthorized, err)
		case constant.ErrPostNotFound:
			return web.MarshalError(c, fiber.StatusNotFound, err)
		default:
			return web.MarshalError(c, fiber.StatusInternalServerError, err)
		}
	}

	c.Status(fiber.StatusNoContent)
	return nil
}
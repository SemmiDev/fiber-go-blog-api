package controllers

import (
	"github.com/SemmiDev/go-blog/internal/app/domain"
	"github.com/SemmiDev/go-blog/internal/helper"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *Server) LikePost(c *fiber.Ctx) error {

	//clear previous error if any
	errList = map[string]string{}

	postID := c.Params("id")

	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	//CHeck if the auth token is valid and  get the user id from it
	metadata, err := s.Tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	//lookup the metadata in redis:
	uid, err := s.Rd.FetchAuth(c.Context(), metadata.TokenUuid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	// check if the user exist:
	user := domain.User{}
	err = s.DB.New().Debug().Model(domain.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// check if the post exist:
	post := domain.Post{}
	err = s.DB.New().Debug().Model(domain.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	like := domain.Like{}
	like.UserID = user.ID
	like.PostID = post.ID

	likeCreated, err := like.SaveLike(s.DB)
	if err != nil {
		formattedError := helper.FormatError(err.Error())
		errList = formattedError
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    likeCreated,
	})
}

func (s *Server) GetLikes(c *fiber.Ctx) error {

	//clear previous error if any
	errList = map[string]string{}

	postID := c.Params("id")

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Check if the post exist:
	post := domain.Post{}
	err = s.DB.Debug().New().Model(domain.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	like := domain.Like{}
	likes, err := like.GetLikesInfo(s.DB, pid)
	if err != nil {
		errList["No_likes"] = "No Likes found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    likes,
	})
}

func (s *Server) UnLikePost(c *fiber.Ctx) error {
	likeID := c.Params("id")

	// Is a valid post id given to us?
	lid, err := strconv.ParseUint(likeID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	//CHeck if the auth token is valid and  get the user id from it
	metadata, err := s.Tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	//lookup the metadata in redis:
	uid, err := s.Rd.FetchAuth(c.Context(), metadata.TokenUuid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	// Check if the post exist
	like := domain.Like{}
	err = s.DB.Debug().New().Model(domain.Like{}).Where("id = ?", lid).Take(&like).Error
	if err != nil {
		errList["No_like"] = "No Like Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Is the authenticated user, the owner of this post?
	if uid != like.UserID {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// If all the conditions are met, delete the post
	_, err = like.DeleteLike(s.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Like deleted",
	})
}

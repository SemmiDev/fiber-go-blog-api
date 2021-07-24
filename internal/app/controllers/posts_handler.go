package controllers

import (
	"github.com/SemmiDev/go-blog/internal/app/domain"
	"github.com/SemmiDev/go-blog/internal/helper"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *Server) CreatePost(c *fiber.Ctx) error {
	errList = map[string]string{}
	post := domain.Post{}

	err := c.BodyParser(&post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	metadata, err := s.Tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	//lookup the metadata in redis:
	userId, err := s.Rd.FetchAuth(c.Context(), metadata.TokenUuid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	// check if the user exist:
	user := domain.User{}
	err = s.DB.Debug().Model(domain.User{}).Where("id = ?", userId).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	post.AuthorID = userId //the authenticated user is the one creating the post
	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	postCreated, err := post.SavePost(s.DB)
	if err != nil {
		errList := helper.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": postCreated,
	})
}

func (s *Server) GetPosts(c *fiber.Ctx) error {

	post := domain.Post{}

	posts, err := post.FindAllPosts(s.DB)
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": posts,
	})
}

func (s *Server) GetPost(c *fiber.Ctx) error {

	postID := c.Params("id")
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	post := domain.Post{}

	postReceived, err := post.FindPostByID(s.DB, pid)
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": postReceived,
	})
}

func (s *Server) UpdatePost(c *fiber.Ctx) error {

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

	//Check if the post exist
	origPost := domain.Post{}
	err = s.DB.New().Debug().Model(domain.Post{}).Where("id = ?", pid).Take(&origPost).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	if uid != origPost.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Start processing the request data
	post := domain.Post{}
	err = c.BodyParser(&post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	post.ID = origPost.ID //this is important to tell the domain the post id to update, the other update field are set above
	post.AuthorID = origPost.AuthorID

	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	postUpdated, err := post.UpdateAPost(s.DB)
	if err != nil {
		errList := helper.FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    postUpdated,
	})
}

func (s *Server) DeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

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
	post := domain.Post{}
	err = s.DB.New().Model(domain.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	if uid != post.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// If all the conditions are met, delete the post
	_, err = post.DeleteAPost(s.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	comment := domain.Comment{}
	like := domain.Like{}

	_, err = comment.DeletePostComments(s.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}
	_, err = like.DeletePostLikes(s.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "post deleted",
	})
}

func (s *Server) GetUserPosts(c *fiber.Ctx) error {
	userID := c.Params("id")
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	post := domain.Post{}
	posts, err := post.FindUserPosts(s.DB, uid)
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": posts,
	})
}

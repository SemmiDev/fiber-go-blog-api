package controllers

import (
	"github.com/SemmiDev/go-blog/internal/app/domain"
	"github.com/SemmiDev/go-blog/internal/helper"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *Server) CreateComment(c *fiber.Ctx) error {
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
	err = s.DB.New().New().Debug().Model(domain.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// check if the post exist:
	post := domain.Post{}
	err = s.DB.New().New().Debug().Model(domain.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	comment := domain.Comment{}
	err = c.BodyParser(&comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// enter the userid and the postid. The comment body is automatically passed
	comment.UserID = uid
	comment.PostID = pid

	comment.Prepare()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}
	commentCreated, err := comment.SaveComment(s.DB)
	if err != nil {
		formattedError := helper.FormatError(err.Error())
		errList = formattedError
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": commentCreated,
	})
}

func (s *Server) GetComments(c *fiber.Ctx) error {
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

	// check if the post exist:
	post := domain.Post{}
	err = s.DB.New().Debug().Model(domain.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No post found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	comment := domain.Comment{}

	comments, err := comment.GetComments(s.DB, pid)
	if err != nil {
		errList["No_comments"] = "No comments found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": comments,
	})
}

func (s *Server) UpdateComment(c *fiber.Ctx) error {

	//clear previous error if any
	errList = map[string]string{}

	commentID := c.Params("id")
	// Check if the post id is valid
	pid, err := strconv.ParseUint(commentID, 10, 64)
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
	origComment := domain.Comment{}
	err = s.DB.New().Debug().Model(domain.Post{}).Where("id = ?", pid).Take(&origComment).Error
	if err != nil {
		errList["No_comment"] = "No Comment Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	if uid != origComment.UserID {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Start processing the request data
	comment := domain.Comment{}
	err = c.BodyParser(&comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	comment.Prepare()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	comment.ID = origComment.ID //this is important to tell the domain the post id to update, the other update field are set above
	comment.UserID = origComment.UserID
	comment.PostID = origComment.PostID

	commentUpdated, err := comment.UpdateAComment(s.DB)
	if err != nil {
		formattedError := helper.FormatError(err.Error())
		errList = formattedError
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": commentUpdated,
	})
}

func (s *Server) DeleteComment(c *fiber.Ctx) error {

	commentID := c.Params("id")
	// Is a valid post id given to us?
	cid, err := strconv.ParseUint(commentID, 10, 64)
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

	// Check if the comment exist
	comment := domain.Comment{}
	err = s.DB.New().Debug().Model(domain.Comment{}).Where("id = ?", cid).Take(&comment).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Is the authenticated user, the owner of this post?
	if uid != comment.UserID {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// If all the conditions are met, delete the post
	_, err = comment.DeleteAComment(s.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": false,
		"message": "Comment deleted",
	})
}

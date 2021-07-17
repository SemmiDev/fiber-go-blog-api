package interfaces

import (
	"fmt"
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	auth2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/auth"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

//Posts struct defines the dependencies that will be used
type Posts struct {
	ui repository2.UserRepository
	ai repository2.PostRepository
	rd auth2.AuthInterface
	tk auth2.TokenInterface
}

//NewPosts constructor
func NewPosts(
	ui repository2.UserRepository,
	ai repository2.PostRepository,
	rd auth2.AuthInterface,
	tk auth2.TokenInterface) *Posts {
	return &Posts{
		ui: ui,
		ai: ai,
		rd: rd,
		tk: tk,
	}
}

func (p *Posts) CreatePost(c *fiber.Ctx) error {
	var post entity2.Post
	var err error

	err = c.BodyParser(&post)
	if err != nil {
		return ERROR(c, fiber.StatusUnprocessableEntity, false, "Cannot parse JSON")
	}

	//check is the user is authenticated first
	var metadata *auth2.AccessDetails
	metadata, err = p.tk.ExtractTokenMetadata(c)
	if err != nil {
		return ERROR(c, fiber.StatusUnauthorized, false, "unauthorized")
	}

	// check if the user exist:
	_, err = p.ui.FindUserByID(metadata.UserId)
	if err != nil {
		return ERROR(c, fiber.StatusUnauthorized, false, err.Error())
	}

	post.AuthorID = metadata.UserId
	post.Prepare()

	var errorMessages map[string]string
	errorMessages = post.Validate()
	if len(errorMessages) > 0 {
		return ERROR(c, fiber.StatusUnprocessableEntity, false, errorMessages)
	}

	var postCreated *entity2.Post
	postCreated, err = p.ai.SavePost(&post)
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err)
	}

	return SUCCESS(c, fiber.StatusCreated, true, postCreated)
}

func (p *Posts) GetPosts(c *fiber.Ctx) error {
	log.Println("LOG")
	posts, err := p.ai.FindAllPosts()
	if err != nil {
		return ERROR(c, fiber.StatusNotFound, false, err)
	}

	return SUCCESS(c, fiber.StatusOK, true, posts)
}

func (p *Posts) GetPost(c *fiber.Ctx) error {
	var err error
	var pid uint64

	pid, err = strconv.ParseUint(c.Params("post_id"), 10, 64)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, err.Error())
	}

	var postReceived *entity2.Post
	postReceived, err = p.ai.FindPostByID(pid)
	if err != nil {
		return ERROR(c, fiber.StatusNotFound, false, err)
	}

	return SUCCESS(c, fiber.StatusOK, true, postReceived)
}

func (p *Posts) GetUserPosts(c *fiber.Ctx) error {
	var err error
	var userID uint64

	userID, err = strconv.ParseUint(c.Params("post_id"), 10, 64)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, err.Error())
	}

	var posts *[]entity2.Post
	posts, err = p.ai.FindUserPosts(userID)
	if err != nil {
		return ERROR(c, fiber.StatusNotFound, false, err)
	}

	return SUCCESS(c, fiber.StatusOK, true, posts)
}

func (p *Posts) DeletePost(c *fiber.Ctx) error {
	pid, err := strconv.ParseUint(c.Params("post_id"), 10, 64)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, err.Error())
	}

	fmt.Println("this is delete post sir")

	//check is the user is authenticated first
	metadata, err := p.tk.ExtractTokenMetadata(c)
	if err != nil {
		return ERROR(c, fiber.StatusUnauthorized, false, "unauthorized")
	}

	post, err := p.ai.FindPostByID(pid)
	if err != nil {
		// Check if the post exist
		return ERROR(c, fiber.StatusNotFound, false, err)
	}

	// Is the authenticated user, the owner of this post?
	if metadata.UserId != post.AuthorID {
		return ERROR(c, fiber.StatusUnauthorized, false, "unauthorized")
	}

	// If all the conditions are met, delete the post
	_, err = p.ai.DeleteAPost(post)
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err)
	}
	//
	//comment := models.Comment{}
	//like := models.Like{}
	//
	//// Also delete the likes and the comments that this post have:
	//_, err = comment.DeletePostComments(server.DB, pid)
	//if err != nil {
	//	errList["Other_error"] = "Please try again later"
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status": http.StatusInternalServerError,
	//		"error":  errList,
	//	})
	//	return
	//}
	//_, err = like.DeletePostLikes(server.DB, pid)
	//if err != nil {
	//	errList["Other_error"] = "Please try again later"
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status": http.StatusInternalServerError,
	//		"error":  errList,
	//	})
	//	return
	//}

	//c.JSON(http.StatusOK, gin.H{
	//	"status":   http.StatusOK,
	//	"response": "Post deleted",
	//})

	return nil
}

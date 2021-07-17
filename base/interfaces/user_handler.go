package interfaces

import (
	"fmt"
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	auth2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/auth"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us repository2.UserRepository
	rd auth2.AuthInterface
	tk auth2.TokenInterface
}

//NewUsers constructor
func NewUsers(us repository2.UserRepository, rd auth2.AuthInterface, tk auth2.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (u *Users) SaveUser(c *fiber.Ctx) error {
	errList := map[string]string{}

	var user entity2.User
	err := c.BodyParser(&user)
	if err != nil {
		return ERROR(c, fiber.StatusUnprocessableEntity, false, "Cannot parse JSON")
	}

	//prepare the request
	user.Prepare()

	//validate the request
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		errList = validateErr
		return ERROR(c, fiber.StatusUnprocessableEntity, false, errList)
	}

	newUser, err := u.us.SaveUser(&user)
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err.Error())
	}
	return SUCCESS(c, fiber.StatusCreated, true, newUser.PublicUser())
}

func (u *Users) FindAllUsers(c *fiber.Ctx) error {
	users := entity2.Users{} //customize user
	var err error
	users, err = u.us.FindAllUsers()
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err.Error())
	}
	return SUCCESS(c, fiber.StatusOK, true, users.PublicUsers())
}

func (u *Users) FindUserByID(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, err.Error())
	}
	user, err := u.us.FindUserByID(userId)
	if err != nil {

		return ERROR(c, fiber.StatusInternalServerError, false, err.Error())
	}
	return SUCCESS(c, fiber.StatusOK, true, user.PublicUser())
}

func (u *Users) UpdateAUser(c *fiber.Ctx) error {
	panic("next time impl y")
}

func (u *Users) UpdateAUserAvatar(c *fiber.Ctx) error {
	panic("next time impl y")
}

func (u *Users) DeleteAUser(c *fiber.Ctx) error {
	//check is the user is authenticated first
	metadata, err := u.tk.ExtractTokenMetadata(c)
	if err != nil {
		return ERROR(c, fiber.StatusUnauthorized, false, "unauthorized")
	}
	userId, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, "invalid request")
	}

	if userId != metadata.UserId {
		return ERROR(c, fiber.StatusUnauthorized, false, "unauthorized")
	}

	_, err = u.us.FindUserByID(metadata.UserId)
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err.Error())
	}

	affected, err := u.us.DeleteAUser(userId)
	if err != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, err.Error())
	}

	return SUCCESS(c, fiber.StatusOK, true, fmt.Sprintf("user with ID %d deleted", affected))
}

func (u *Users) UpdatePassword(c *fiber.Ctx) error {
	panic("")
}

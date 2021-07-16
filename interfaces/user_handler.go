package interfaces

import (
	"github.com/SemmiDev/fiber-go-blog/application"
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
	"github.com/SemmiDev/fiber-go-blog/infrastructure/auth"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

//NewUsers constructor
func NewUsers(us application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (u *Users) SaveUser(c *fiber.Ctx) error {
	var user entity.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	user.Prepare()

	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": validateErr,
		})
	}

	newUser, err := u.us.SaveUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newUser.PublicUser(),
	})
}

func (u *Users) FindAllUsers(c *fiber.Ctx) error {
	users := entity.Users{} //customize user
	var err error
	users, err = u.us.FindAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    users.PublicUsers(),
	})
}

func (u *Users) FindUserByID(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	user, err := u.us.FindUserByID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    user.PublicUser(),
	})
}

func (u *Users) UpdateAUser(c *fiber.Ctx) error {
	panic("next time impl y")
}

func (u *Users) UpdateAUserAvatar(c *fiber.Ctx) error {
	panic("next time impl y")
}

func (u *Users) DeleteAUser(c *fiber.Ctx) error {
	panic("next time impl y")
}

func (u *Users) UpdatePassword(c *fiber.Ctx) error {
	panic("next time impl y")
}

package interfaces

import (
	"fmt"
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	auth2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
)

//Authenticate struct defines the dependencies that will be used
type Authenticate struct {
	us repository2.UserRepository
	rd auth2.AuthInterface
	tk auth2.TokenInterface
}

//NewAuthenticate constructor
func NewAuthenticate(us repository2.UserRepository, rd auth2.AuthInterface, tk auth2.TokenInterface) *Authenticate {
	return &Authenticate{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (au *Authenticate) Login(c *fiber.Ctx) error {
	var user *entity2.User

	err := c.BodyParser(&user)
	if err != nil {
		return ERROR(c, fiber.StatusBadRequest, false, "Cannot parse JSON")
	}

	user.Prepare()
	validateUser := user.Validate("LOGIN")
	if len(validateUser) > 0 {
		//validate request:
		return ERROR(c, fiber.StatusUnprocessableEntity, false, validateUser)
	}
	u, userErr := au.us.GetUserByEmailAndPassword(user)
	if userErr != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, userErr)
	}
	ts, tErr := au.tk.CreateToken(u.ID)
	if tErr != nil {
		return ERROR(c, fiber.StatusUnprocessableEntity, false, tErr.Error())
	}
	saveErr := au.rd.CreateAuth(c.Context(), u.ID, ts)
	if saveErr != nil {
		return ERROR(c, fiber.StatusInternalServerError, false, saveErr.Error())
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["id"] = u.ID
	userData["email"] = u.Email
	userData["avatar_path"] = u.AvatarPath
	userData["username"] = u.Username

	return SUCCESS(c, fiber.StatusOK, true, userData)
}

func (au *Authenticate) Logout(c *fiber.Ctx) error {
	//check is the user is authenticated first
	metadata, err := au.tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	//if the access token exist and it is still valid, then delete both the access token and the refresh token
	deleteErr := au.rd.DeleteTokens(c.Context(), metadata)
	if deleteErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": deleteErr.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged out",
	})
}

//Refresh is the function that uses the refresh_token to generate new pairs of refresh and access tokens.
func (au *Authenticate) Refresh(c *fiber.Ctx) error {
	mapToken := map[string]string{}
	err := c.BodyParser(&mapToken)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//any error may be due to token expiration
	if err != nil {
		return ERROR(c, fiber.StatusUnauthorized, false, err.Error())
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return ERROR(c, fiber.StatusUnauthorized, false, "Unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			//Since token is valid, get the uuid:
			return ERROR(c, fiber.StatusUnprocessableEntity, false, "Cannot get uuid")
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return ERROR(c, fiber.StatusUnprocessableEntity, false, "Error occurred")
		}
		delErr := au.rd.DeleteRefresh(c.Context(), refreshUuid)
		if delErr != nil { //if any goes wrong
			//Delete the previous Refresh Token
			return ERROR(c, fiber.StatusUnauthorized, false, "Unauthorized")
		}
		ts, createErr := au.tk.CreateToken(userId)
		if createErr != nil {
			//Create new pairs of refresh and access tokens
			return ERROR(c, fiber.StatusForbidden, false, createErr.Error())
		}
		saveErr := au.rd.CreateAuth(c.Context(), userId, ts)
		if saveErr != nil {
			//save the tokens metadata to redis
			return ERROR(c, fiber.StatusForbidden, false, saveErr.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return SUCCESS(c, fiber.StatusCreated, true, tokens)
	} else {
		return SUCCESS(c, fiber.StatusUnauthorized, false, "Refresh token expired")
	}
}

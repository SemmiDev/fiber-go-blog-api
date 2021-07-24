package controllers

import (
	"errors"
	"fmt"
	"github.com/SemmiDev/go-blog/internal/app/domain"
	"github.com/SemmiDev/go-blog/internal/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
)

func (s *Server) Login(c *fiber.Ctx) error {
	var user domain.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}

	user.Prepare()

	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": validateUser,
		})
	}

	u, userErr := GetUserByEmailAndPassword(s.DB, &user)
	if userErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": userErr.Error(),
		})
	}

	ts, tErr := s.Tk.CreateToken(u.ID)
	if tErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": tErr.Error(),
		})
	}

	saveErr := s.Rd.CreateAuth(c.Context(), u.ID, ts)
	if saveErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": saveErr.Error(),
		})
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["id"] = u.ID
	userData["username"] = u.Username
	userData["name"] = u.Name

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userData,
	})
}

func (s *Server) Logout(c *fiber.Ctx) error {
	//check is the user is authenticated first
	metadata, err := s.Tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	//if the access token exist and it is still valid, then delete both the access token and the refresh token
	deleteErr := s.Rd.DeleteTokens(c.Context(), metadata)
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
func (s *Server) Refresh(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "Cannot get uuid",
			})
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "Error occurred",
			})
		}
		//Delete the previous Refresh Token
		delErr := s.Rd.DeleteRefresh(c.Context(), refreshUuid)
		if delErr != nil { //if any goes wrong
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := s.Tk.CreateToken(userId)
		if createErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": createErr.Error(),
			})
		}
		//save the tokens metadata to redis
		saveErr := s.Rd.CreateAuth(c.Context(), userId, ts)
		if saveErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": saveErr.Error(),
			})
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    tokens,
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": true,
			"message": "Refresh token expired",
		})
	}
}

func GetUserByEmailAndPassword(db *gorm.DB, u *domain.User) (*domain.User, error) {
	var user domain.User
	err := db.New().Debug().Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}

	// verify
	if err := helper.CheckPassword(u.Password, user.Password); err != nil {
		return nil, errors.New(err.Error())
	}

	return &user, nil
}

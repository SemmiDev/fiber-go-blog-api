package app

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

func (s *Server) CreateUser(c *fiber.Ctx) error {
	errList = map[string]string{}
	var user User

	err := c.BodyParser(&user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	errorMessages := user.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	userCreated, err := user.SaveUser(s.DB)
	if err != nil {
		formattedError := FormatError(err.Error())
		errList = formattedError
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    userCreated,
	})
}

func (s *Server) GetUsers(c *fiber.Ctx) error {
	//clear previous error if any
	errList = map[string]string{}
	user := User{}

	users, err := user.FindAllUsers(s.DB)
	if err != nil {
		errList["No_user"] = "No User Found"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": true,
			"data":    errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    users,
	})
}

func (s *Server) GetUser(c *fiber.Ctx) error {
	errList = map[string]string{}
	userID := c.Params("id")

	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	user := User{}
	userGotten, err := user.FindUserByID(s.DB, uid)
	if err != nil {
		errList["No_user"] = "No User Found"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    userGotten,
	})
}

func (s *Server) UpdateUser(c *fiber.Ctx) error {
	errList = map[string]string{}

	userID := c.Params("id")
	// Check if the user id is valid
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	//Check if the auth token is valid and  get the user id from it
	metadata, err := s.Tk.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}

	log.Println(metadata.TokenUuid)

	//lookup the metadata in redis:
	tokenID, err := s.Rd.FetchAuth(c.Context(), metadata.TokenUuid)
	if err != nil {
		log.Println("yes di redis salah e")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "unauthorized",
		})
	}
	log.Println("after redis")

	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != uid {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Start processing the request
	requestBody := map[string]string{}
	err = c.BodyParser(&requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Check for previous details
	formerUser := User{}
	err = s.DB.New().Debug().Model(User{}).Where("id = ?", uid).Take(&formerUser).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	newUser := User{}

	//When current password has content.
	if requestBody["current_password"] == "" && requestBody["new_password"] != "" {
		errList["Empty_current"] = "Please Provide current password"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] == "" {
		errList["Empty_new"] = "Please Provide new password"
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	if requestBody["current_password"] != "" && requestBody["new_password"] != "" {
		//Also check if the new password
		if len(requestBody["new_password"]) < 6 {
			errList["Invalid_password"] = "Password should be at least 6 characters"
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": errList,
			})
		}

		//if they do, check that the former password is correct
		err = CheckPassword(requestBody["current_password"], formerUser.Password)
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			errList["Password_mismatch"] = "The password not correct"
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": errList,
			})
		}

		//update both the password and the email
		newUser.Username = formerUser.Username //remember, you cannot update the username
		newUser.Email = requestBody["email"]

		hashed, _ := HashPassword(requestBody["new_password"])
		newUser.Password = hashed
	} else {
		newUser.Password = formerUser.Password
	}

	//The password fields not entered, so update only the email
	newUser.Username = formerUser.Username

	newUser.Prepare()
	errorMessages := newUser.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	log.Println("----------------------------- JAUH GAN")
	log.Println(newUser)
	log.Println("-----------------------------")

	updatedUser, err := newUser.UpdateAUser(s.DB, uid)
	if err != nil {
		errList := FormatError(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    updatedUser,
	})
}

func (s *Server) DeleteUser(c *fiber.Ctx) error {
	errList = map[string]string{}
	userID := c.Params("id")

	uid, err := strconv.ParseUint(userID, 10, 64)
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
			"message": "Unauthorized",
		})
	}

	// If the id is not the authenticated user id
	if metadata.UserId != 0 && metadata.UserId != uid {
		errList["Unauthorized"] = "Unauthorized"
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	err = s.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	user := User{}
	_, err = user.DeleteAUser(s.DB, uid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	// Also delete the posts, likes and the comments that this user created if any:
	comment := Comment{}
	like := Like{}
	post := Post{}

	_, err = post.DeleteUserPosts(s.DB, uid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	_, err = comment.DeleteUserComments(s.DB, uid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	_, err = like.DeleteUserLikes(s.DB, uid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": errList,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    "user deleted",
	})
}

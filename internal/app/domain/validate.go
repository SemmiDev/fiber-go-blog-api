package domain

import (
	"errors"
	"github.com/badoux/checkmail"
	"strings"
)

// Validate for user
func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}

	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		if u.Name == "" {
			err = errors.New("Required Name")
			errorMessages["Required_name"] = err.Error()
		}
		if u.Username == "" {
			err = errors.New("Required Username")
			errorMessages["Required_username"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

// Validate for post
func (p *Post) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()

	}
	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

// Validate for comment
func (c *Comment) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMessages["Required_body"] = err.Error()
		}
	default:
		if c.Body == "" {
			err = errors.New("Required Comment")
			errorMessages["Required_body"] = err.Error()
		}
	}
	return errorMessages
}

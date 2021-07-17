package entity

import (
	security2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/security"
	"github.com/badoux/checkmail"
	"html"
	"os"
	"strings"
	"time"
)

type User struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	FullName   string    `gorm:"size:255;not null" json:"name"`
	Username   string    `gorm:"size:255;not null;unique" json:"username"`
	Email      string    `gorm:"size:100;not null;unique" json:"email"`
	Password   string    `gorm:"size:100;not null;" json:"password"`
	AvatarPath string    `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicUser struct {
	ID       uint64 `gorm:"primary_key;auto_increment" json:"id"`
	FullName string `gorm:"size:100;not null;" json:"name"`
}

type Users []*User

// PublicUsers So that we dont expose the user's email address and password to the world
func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

// PublicUser So that we dont expose the user's email address and password to the world
func (u *User) PublicUser() interface{} {
	return &PublicUser{
		ID:       u.ID,
		FullName: u.FullName,
	}
}

func (u *User) BeforeSave() error {
	hashedPassword, err := security2.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.FullName = html.EscapeString(strings.TrimSpace(u.FullName))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.FullName = strings.Title(strings.ToLower(u.FullName))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) AfterFind() (err error) {
	if u.AvatarPath != "" {
		u.AvatarPath = os.Getenv("DO_SPACES_URL") + u.AvatarPath
	}
	return nil
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch action {
	case "UPDATE":
		if u.Email == "" {
			errorMessages["Required_email"] = ErrRequiredEmail.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["Invalid_email"] = ErrInvalidEmail.Error()
			}
		}
	case "LOGIN":
		if u.Password == "" {
			errorMessages["Required_password"] = ErrRequiredPassword.Error()
		}
		if u.Email == "" {
			errorMessages["Required_email"] = ErrRequiredEmail.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["Invalid_email"] = ErrInvalidEmail.Error()
			}
		}
	case "FORGOT_PASSWORD":
		if u.Email == "" {
			errorMessages["Required_email"] = ErrRequiredEmail.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["Invalid_email"] = ErrInvalidEmail.Error()
			}
		}
	default:
		if u.Username == "" {
			errorMessages["Required_username"] = ErrRequiredUsername.Error()
		}
		if u.Password == "" {
			errorMessages["Required_password"] = ErrRequiredPassword.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			errorMessages["Invalid_password"] = ErrEmailRule.Error()
		}
		if u.Email == "" {
			errorMessages["Required_email"] = ErrRequiredEmail.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["Invalid_email"] = ErrInvalidEmail.Error()
			}
		}
	}
	return errorMessages
}

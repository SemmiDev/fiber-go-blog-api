package app

import (
	"html"
	"log"
	"strings"
	"time"
)

type User struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PublicUser struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (u *User) PublicUser() *PublicUser {
	return &PublicUser{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
	}
}

func (u *User) BeforeSave() error {
	hashedPassword, err := HashPassword(u.Password)
	log.Println("--------------------------------")
	log.Println(u.Password)
	log.Println(hashedPassword)
	log.Println("---------------------------------")

	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) Prepare() {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

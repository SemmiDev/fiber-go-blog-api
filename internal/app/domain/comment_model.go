package domain

import (
	"html"
	"strings"
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	Body      string    `gorm:"text;not null;" json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

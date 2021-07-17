package entity

import (
	"html"
	"strings"
	"time"
)

type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint64    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	//p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate() map[string]string {
	var errorMessages = make(map[string]string)

	if p.Title == "" {
		errorMessages["Required_title"] = ErrRequiredTitle.Error()
	}
	if p.Content == "" {
		errorMessages["Required_content"] = ErrRequiredContent.Error()
	}
	if p.AuthorID < 1 {
		errorMessages["Required_author"] = ErrRequiredAuthor.Error()
	}
	return errorMessages
}

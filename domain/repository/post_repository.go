package repository

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
	"github.com/jinzhu/gorm"
)

type PostRepository interface {
	SavePost(db *gorm.DB) (*entity.Post, error)
	FindAllPosts(db *gorm.DB) (*[]entity.Post, error)
	FindPostByID(pid uint64) (*entity.Post, error)
	UpdateAPost() (*entity.Post, error)
	DeleteAPost() (int64, error)
	FindUserPosts(uid uint64) (*[]entity.Post, error)
	DeleteUserPosts(uid uint64) (int64, error)
}

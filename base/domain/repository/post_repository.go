package repository

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
)

type PostRepository interface {
	SavePost(p *entity2.Post) (*entity2.Post, error)
	FindAllPosts() (*[]entity2.Post, error)
	FindPostByID(pid uint64) (*entity2.Post, error)
	UpdateAPost(p *entity2.Post) (*entity2.Post, error)
	DeleteAPost(p *entity2.Post) (int64, error)
	FindUserPosts(uid uint64) (*[]entity2.Post, error)
	DeleteUserPosts(uid uint64) (int64, error)
}

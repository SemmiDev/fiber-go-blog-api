package repository

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
)

type CommentRepository interface {
	SaveComment(*entity.Comment, error)
	GetComments(pid uint64) (*[]entity.Comment, error)
	UpdateAComment(*entity.Comment, error)
	DeleteAComment(int64, error)
	DeleteUserComments(uid uint64) (int64, error)
	DeletePostComments(pid uint64) (int64, error)
}

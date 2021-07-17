package repository

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
)

type CommentRepository interface {
	SaveComment(*entity2.Comment, error)
	GetComments(pid uint64) (*[]entity2.Comment, error)
	UpdateAComment(*entity2.Comment, error)
	DeleteAComment(int64, error)
	DeleteUserComments(uid uint64) (int64, error)
	DeletePostComments(pid uint64) (int64, error)
}

package repository

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
)

type LikeRepository interface {
	SaveLike() (*entity2.Like, error)
	DeleteLike() (*entity2.Like, error)
	GetLikesInfo(pid uint64) (*[]entity2.Like, error)
	DeleteUserLikes(uid uint64) (int64, error)
	DeletePostLikes(pid uint64) (int64, error)
}

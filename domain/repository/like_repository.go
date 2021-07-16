package repository

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
)

type LikeRepository interface {
	SaveLike() (*entity.Like, error)
	DeleteLike() (*entity.Like, error)
	GetLikesInfo(pid uint64) (*[]entity.Like, error)
	DeleteUserLikes(uid uint64) (int64, error)
	DeletePostLikes(pid uint64) (int64, error)
}

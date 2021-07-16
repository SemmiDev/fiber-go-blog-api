package repository

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
)

type UserRepository interface {
	SaveUser(user *entity.User) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
	FindUserByID(uid uint64) (*entity.User, error)
	UpdateAUser(user *entity.User, uid int64) (*entity.User, error)
	UpdateAUserAvatar(user *entity.User, uid int64) (*entity.User, error)
	DeleteAUser(uid int64) (int64, error)
	UpdatePassword(user *entity.User) error
}

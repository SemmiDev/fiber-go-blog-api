package repository

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
)

type UserRepository interface {
	SaveUser(user *entity2.User) (*entity2.User, error)
	FindAllUsers() ([]*entity2.User, error)
	FindUserByID(uid uint64) (*entity2.User, error)
	UpdateAUser(user *entity2.User, uid uint64) (*entity2.User, error)
	UpdateAUserAvatar(user *entity2.User, uid uint64) (*entity2.User, error)
	DeleteAUser(uid uint64) (uint64, error)
	UpdatePassword(user *entity2.User) error
	GetUserByEmailAndPassword(*entity2.User) (*entity2.User, error)
}

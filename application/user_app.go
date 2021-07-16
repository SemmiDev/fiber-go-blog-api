package application

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
	"github.com/SemmiDev/fiber-go-blog/domain/repository"
)

type userApp struct {
	us repository.UserRepository
}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, error) {
	return u.us.SaveUser(user)
}

func (u *userApp) FindAllUsers() ([]*entity.User, error) {
	return u.us.FindAllUsers()
}

func (u *userApp) FindUserByID(uid uint64) (*entity.User, error) {
	return u.us.FindUserByID(uid)
}

func (u *userApp) UpdateAUser(user *entity.User, uid int64) (*entity.User, error) {
	return u.us.UpdateAUser(user, uid)
}

func (u *userApp) UpdateAUserAvatar(user *entity.User, uid int64) (*entity.User, error) {
	return u.us.UpdateAUserAvatar(user, uid)
}

func (u *userApp) DeleteAUser(uid int64) (int64, error) {
	return u.us.DeleteAUser(uid)
}

func (u *userApp) UpdatePassword(user *entity.User) error {
	return u.us.UpdatePassword(user)
}

var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(user *entity.User) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
	FindUserByID(uid uint64) (*entity.User, error)
	UpdateAUser(user *entity.User, uid int64) (*entity.User, error)
	UpdateAUserAvatar(user *entity.User, uid int64) (*entity.User, error)
	DeleteAUser(uid int64) (int64, error)
	UpdatePassword(user *entity.User) error
}

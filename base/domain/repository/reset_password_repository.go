package repository

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
)

type ResetPasswordRepository interface {
	SaveDetails() (*entity2.ResetPassword, error)
	DeleteDetails() (int64, error)
}

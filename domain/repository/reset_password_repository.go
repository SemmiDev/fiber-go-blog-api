package repository

import (
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
)

type ResetPasswordRepository interface {
	SaveDetails() (*entity.ResetPassword, error)
	DeleteDetails() (int64, error)
}

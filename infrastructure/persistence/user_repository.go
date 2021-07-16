package persistence

import (
	"errors"
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
	"github.com/SemmiDev/fiber-go-blog/domain/repository"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func (u *UserRepo) SaveUser(user *entity.User) (*entity.User, error) {
	var err error
	err = u.db.Debug().Create(user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

func (u *UserRepo) FindAllUsers() ([]*entity.User, error) {
	var err error
	var users []*entity.User
	err = u.db.Debug().Model(&entity.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return []*entity.User{}, err
	}
	return users, err
}

func (u *UserRepo) FindUserByID(uid uint64) (*entity.User, error) {
	var err error
	var user entity.User

	err = u.db.Debug().Model(entity.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity.User{}, errors.New("User Not Found")
	}
	return &user, err
}

func (u *UserRepo) UpdateAUser(user *entity.User, uid int64) (*entity.User, error) {
	if user.Password != "" {
		// To hash the password
		err := user.BeforeSave()
		if err != nil {
			log.Fatal(err)
		}

		u.db = u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).UpdateColumns(
			map[string]interface{}{
				"password":  user.Password,
				"email":     user.Email,
				"update_at": time.Now(),
			},
		)
	}

	u.db = u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).UpdateColumns(
		map[string]interface{}{
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)
	if u.db.Error != nil {
		return &entity.User{}, u.db.Error
	}

	// This is the display the updated user
	err := u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil

}

func (u *UserRepo) UpdateAUserAvatar(user *entity.User, uid int64) (*entity.User, error) {
	u.db = u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).UpdateColumns(
		map[string]interface{}{
			"avatar_path": user.AvatarPath,
			"update_at":   time.Now(),
		},
	)
	if u.db.Error != nil {
		return &entity.User{}, u.db.Error
	}
	// This is the display the updated user
	err := u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

func (u *UserRepo) DeleteAUser(uid int64) (int64, error) {
	u.db = u.db.Debug().Model(&entity.User{}).Where("id = ?", uid).Take(&entity.User{}).Delete(&entity.User{})

	if u.db.Error != nil {
		return 0, u.db.Error
	}
	return u.db.RowsAffected, nil
}

func (u *UserRepo) UpdatePassword(user *entity.User) error {
	// To hash the password
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	u.db = u.db.Debug().Model(&entity.User{}).Where("email = ?", user.Email).Take(&entity.User{}).UpdateColumns(
		map[string]interface{}{
			"password":  user.Password,
			"update_at": time.Now(),
		},
	)
	if u.db.Error != nil {
		return u.db.Error
	}
	return nil
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

//UserRepo implements the repository.UserRepository interface
var _ repository.UserRepository = &UserRepo{}

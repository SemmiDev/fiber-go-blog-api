package persistence

import (
	"errors"
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	security2 "github.com/SemmiDev/fiber-go-blog/base/infrastructure/security"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func (u *UserRepo) GetUserByEmailAndPassword(user *entity2.User) (*entity2.User, error) {
	var me entity2.User
	err := u.db.Debug().Where("email = ?", user.Email).Take(&me).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("database error")
	}
	//Verify the password
	err = security2.Verify(user.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("incorrect password")
	}
	return &me, nil
}

func (u *UserRepo) SaveUser(user *entity2.User) (*entity2.User, error) {
	var err error
	err = u.db.Debug().Create(user).Error
	if err != nil {
		return &entity2.User{}, err
	}
	return user, nil
}

func (u *UserRepo) FindAllUsers() ([]*entity2.User, error) {
	var err error
	var users []*entity2.User
	err = u.db.Debug().Model(&entity2.User{}).Limit(100).Find(&users).Error
	if err != nil {
		return []*entity2.User{}, err
	}
	return users, err
}

func (u *UserRepo) FindUserByID(uid uint64) (*entity2.User, error) {
	var err error
	var user entity2.User

	err = u.db.Debug().Model(entity2.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		return &entity2.User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &entity2.User{}, errors.New("User Not Found")
	}
	return &user, err
}

func (u *UserRepo) UpdateAUser(user *entity2.User, uid uint64) (*entity2.User, error) {
	if user.Password != "" {
		// To hash the password
		err := user.BeforeSave()
		if err != nil {
			log.Fatal(err)
		}

		u.db = u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(&entity2.User{}).UpdateColumns(
			map[string]interface{}{
				"password":  user.Password,
				"email":     user.Email,
				"update_at": time.Now(),
			},
		)
	}

	u.db = u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(&entity2.User{}).UpdateColumns(
		map[string]interface{}{
			"email":     user.Email,
			"update_at": time.Now(),
		},
	)
	if u.db.Error != nil {
		return &entity2.User{}, u.db.Error
	}

	// This is the display the updated user
	err := u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(user).Error
	if err != nil {
		return &entity2.User{}, err
	}
	return user, nil

}

func (u *UserRepo) UpdateAUserAvatar(user *entity2.User, uid uint64) (*entity2.User, error) {
	u.db = u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(&entity2.User{}).UpdateColumns(
		map[string]interface{}{
			"avatar_path": user.AvatarPath,
			"update_at":   time.Now(),
		},
	)
	if u.db.Error != nil {
		return &entity2.User{}, u.db.Error
	}
	// This is the display the updated user
	err := u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &entity2.User{}, err
	}
	return user, nil
}

func (u *UserRepo) DeleteAUser(uid uint64) (uint64, error) {
	u.db = u.db.Debug().Model(&entity2.User{}).Where("id = ?", uid).Take(&entity2.User{}).Delete(&entity2.User{})

	if u.db.Error != nil {
		return 0, u.db.Error
	}

	log.Println(u.db.RowsAffected)

	return uint64(u.db.RowsAffected), nil
}

func (u *UserRepo) UpdatePassword(user *entity2.User) error {
	// To hash the password
	err := user.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	u.db = u.db.Debug().Model(&entity2.User{}).Where("email = ?", user.Email).Take(&entity2.User{}).UpdateColumns(
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
var _ repository2.UserRepository = &UserRepo{}

package domain

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	log.Println(u)
	err = db.Debug().Create(u).Error
	log.Println(u)
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	var users []User
	err = db.New().Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint64) (*User, error) {
	var err error
	err = db.New().Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint64) (*User, error) {
	db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"updated_at": time.Now(),
		},
	)

	if u.Password != "" {
		db = db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
			map[string]interface{}{
				"password": u.Password,
				"email":    u.Email,
			},
		)
		if db.Error != nil {
			return &User{}, db.Error
		}

		err := db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
		if err != nil {
			return &User{}, err
		}

		return u, nil
	} else {
		db = db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
			map[string]interface{}{
				"email": u.Email,
			},
		)
		if db.Error != nil {
			return &User{}, db.Error
		}

		err := db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
		if err != nil {
			return &User{}, err
		}

		return u, nil
	}
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint64) (int64, error) {
	db = db.New().Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) UpdatePassword(db *gorm.DB) error {
	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.New().Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

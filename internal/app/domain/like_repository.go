package domain

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {
	// Check if the auth user has liked this post before:
	err := db.New().Debug().Model(&Like{}).Where("post_id = ? AND user_id = ?", l.PostID, l.UserID).Take(&l).Error
	if err != nil {
		if err.Error() == "record not found" {
			// The user has not liked this post before, so lets save incomming like:
			err = db.New().Debug().Model(&Like{}).Create(&l).Error
			if err != nil {
				return &Like{}, err
			}
		}
	} else {
		// The user has liked it before, so create a custom error message
		err = errors.New("double like")
		return &Like{}, err
	}
	return l, nil
}

func (l *Like) DeleteLike(db *gorm.DB) (*Like, error) {
	var err error
	var deletedLike *Like

	err = db.New().Debug().Model(Like{}).Where("id = ?", l.ID).Take(&l).Error
	if err != nil {
		return &Like{}, err
	} else {
		//If the like exist, save it in deleted like and delete it
		deletedLike = l
		db = db.New().Debug().Model(&Like{}).Where("id = ?", l.ID).Take(&Like{}).Delete(&Like{})
		if db.Error != nil {
			fmt.Println("cant delete like: ", db.Error)
			return &Like{}, db.Error
		}
	}
	return deletedLike, nil
}

func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error) {

	var likes []Like
	err := db.New().Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Error
	if err != nil {
		return &[]Like{}, err
	}
	return &likes, err
}

func (l *Like) DeleteUserLikes(db *gorm.DB, uid uint64) (int64, error) {
	var likes []Like
	db = db.New().Debug().Model(&Like{}).Where("user_id = ?", uid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	var likes []Like
	db = db.New().Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

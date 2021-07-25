package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	err := db.New().Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.New().Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) GetComments(db *gorm.DB, pid uint64) (*[]Comment, error) {

	var comments []Comment
	err := db.New().Debug().Model(&Comment{}).Where("post_id = ?", pid).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i := range comments {
			err := db.New().Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].User).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, err
}

func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {
	var err error

	err = db.New().Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Body: c.Body, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}

	fmt.Println("this is the comment body: ", c.Body)
	if c.ID != 0 {
		err = db.New().Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {

	db = db.New().Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteUserComments When a user is deleted, we also delete the comments that the user had
func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint64) (int64, error) {
	var comments []Comment
	db = db.New().Debug().Model(&Comment{}).Where("user_id = ?", uid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeletePostComments When a post is deleted, we also delete the comments that the post had
func (c *Comment) DeletePostComments(db *gorm.DB, pid uint64) (int64, error) {
	var comments []Comment
	db = db.New().Debug().Model(&Comment{}).Where("post_id = ?", pid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

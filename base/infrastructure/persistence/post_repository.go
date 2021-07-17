package persistence

import (
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type PostRepo struct {
	db       *gorm.DB
	userRepo *UserRepo
}

func (pr *PostRepo) SavePost(p *entity2.Post) (*entity2.Post, error) {
	log.Println("PERSISTENCE PASSED")
	var err error
	err = pr.db.Debug().Model(&entity2.Post{}).Create(&p).Error
	if err != nil {
		return &entity2.Post{}, err
	}
	if p.ID != 0 {
		err = pr.db.Debug().Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &entity2.Post{}, err
		}
	}
	return p, nil
}

func (pr *PostRepo) FindAllPosts() (*[]entity2.Post, error) {
	var err error
	var posts []entity2.Post
	err = pr.db.Debug().Model(&entity2.Post{}).Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return &[]entity2.Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := pr.db.Debug().Model(&entity2.User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]entity2.Post{}, err
			}
		}
	}
	return &posts, nil
}

func (pr *PostRepo) FindPostByID(pid uint64) (*entity2.Post, error) {
	var p entity2.Post
	var err error

	err = pr.db.Debug().Model(&entity2.Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &entity2.Post{}, err
	}
	if p.ID != 0 {
		err = pr.db.Debug().Model(&entity2.User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &entity2.Post{}, err
		}
	}
	return &p, nil
}

func (pr *PostRepo) UpdateAPost(p *entity2.Post) (*entity2.Post, error) {
	var err error

	err = pr.db.Debug().Model(&entity2.Post{}).Where("id = ?", p.ID).Updates(
		entity2.Post{
			Title:     p.Title,
			Content:   p.Content,
			UpdatedAt: time.Now(),
		}).Error

	if err != nil {
		return &entity2.Post{}, err
	}
	if p.ID != 0 {
		err = pr.db.Debug().Model(&entity2.User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &entity2.Post{}, err
		}
	}
	return p, nil
}

func (pr *PostRepo) DeleteAPost(p *entity2.Post) (int64, error) {
	pr.db = pr.db.Debug().Model(&entity2.Post{}).Where("id = ?", p.ID).Take(&entity2.Post{}).Delete(&entity2.Post{})
	if pr.db.Error != nil {
		return 0, pr.db.Error
	}
	return pr.db.RowsAffected, nil
}

func (pr *PostRepo) FindUserPosts(uid uint64) (*[]entity2.Post, error) {
	var err error
	var posts []entity2.Post
	err = pr.db.Debug().Model(&entity2.Post{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return &[]entity2.Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := pr.db.Debug().Model(&entity2.User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]entity2.Post{}, err
			}
		}
	}
	return &posts, nil
}

func (pr *PostRepo) DeleteUserPosts(uid uint64) (int64, error) {
	var posts []entity2.Post
	pr.db = pr.db.Debug().Model(&entity2.Post{}).Where("author_id = ?", uid).Find(&posts).Delete(&posts)
	if pr.db.Error != nil {
		return 0, pr.db.Error
	}
	return pr.db.RowsAffected, nil
}

func NewPostRepository(db *gorm.DB, repo *UserRepo) *PostRepo {
	return &PostRepo{db, repo}
}

//PostRepo implements the repository.PostRepository interface
var _ repository2.PostRepository = &PostRepo{}

package persistence

import (
	"fmt"
	"github.com/SemmiDev/fiber-go-blog/domain/entity"
	"github.com/SemmiDev/fiber-go-blog/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repositories struct {
	User          repository.UserRepository
	Post          repository.PostRepository
	Comment       repository.CommentRepository
	Like          repository.LikeRepository
	ResetPassword repository.ResetPasswordRepository
	db            *gorm.DB
}

func NewRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		//Post:          NewPostRepository(db),
		//Comment:       NewCommentRepository(db),
		//Like:          NewLikeRepository(db),
		//ResetPassword: NewResetPasswordRepository(db),
		db: db,
	}, nil
}

// Close closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

// Automigrate This migrate all tables
func (s *Repositories) Automigrate() error {
	return s.db.AutoMigrate(
		&entity.User{},
		&entity.Post{},
		&entity.Comment{},
		&entity.Like{},
		&entity.ResetPassword{},
	).Error
}

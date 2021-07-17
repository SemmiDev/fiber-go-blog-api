package persistence

import (
	"fmt"
	entity2 "github.com/SemmiDev/fiber-go-blog/base/domain/entity"
	repository2 "github.com/SemmiDev/fiber-go-blog/base/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type Repositories struct {
	User          repository2.UserRepository
	Post          repository2.PostRepository
	Comment       repository2.CommentRepository
	Like          repository2.LikeRepository
	ResetPassword repository2.ResetPasswordRepository
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
		&entity2.User{},
		&entity2.Post{},
		&entity2.Comment{},
		&entity2.Like{},
		&entity2.ResetPassword{},
	).Error
}

// DropTables This drop all tables (for dev only)
func (s *Repositories) DropTables() {
	err := s.db.DropTableIfExists(
		&entity2.User{},
		&entity2.Post{},
		&entity2.Comment{},
		&entity2.Like{},
		&entity2.ResetPassword{},
	).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
}

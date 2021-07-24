package seeder

import (
	"github.com/SemmiDev/go-blog/internal/app/domain"
	"github.com/jinzhu/gorm"
	"log"
)

//var hashed1, _ = helper.HashPassword("sammidev")
//var hashed2, _ = helper.HashPassword("izzah")
//var hashed3, _ = helper.HashPassword("sammidevizzah")

var users = []domain.User{
	{
		Name:     "sammi",
		Username: "sammi",
		Email:    "sammi@gmail.com",
		Password: "sammi",
	},
	{
		Name:     "izzah",
		Username: "izzah",
		Email:    "izzah@gmail.com",
		Password: "sammi",
	},
	{
		Name:     "sammiizzah",
		Username: "sammiizzah",
		Email:    "sammiizzah@gmail.com",
		Password: "sammi",
	},
}

var posts = []domain.Post{
	{
		Title:   "HMMM",
		Content: "Aku menaruhmu terlalu dalam di hati, sehingga untuk menghapusmu, aku seperti menyakiti diri sendiri uakh dei!",
	},
	{
		Title:   "I LOVE U",
		Content: "ahabbakalladzi ahbabtani lahu",
	},
	{
		Title:   "RINTIK",
		Content: "Meski mentari sudah menenggelamkan diri, awan kelabu sudah tak nampak lagi, tapi aku masih tetap menanti rintik hujan membasahi bumi",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().RemoveForeignKey("author_id", "users(id)").Error
	if err != nil {
		log.Fatalf("cannot remove foreign key: %v", err)
	}

	err = db.Debug().DropTableIfExists(
		&domain.Post{},
		&domain.User{},
		&domain.Like{},
		&domain.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(
		&domain.Post{},
		&domain.User{},
		&domain.Like{},
		&domain.Comment{}).Error

	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&domain.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&domain.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&domain.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}

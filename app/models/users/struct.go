package users

import (
	"database/sql"
	"time"

	"github.com/daison12006013/gorvel/databases"
	"github.com/daison12006013/gorvel/pkg/errors"
	"github.com/daison12006013/gorvel/pkg/paginate"
	"gorm.io/gorm"
)

const Table = "users"
const PrimaryKey = "id"

type Paginate struct {
	paginate.Paginate
}

type Attributes struct {
	// gorm.Model
	ID              uint           `gorm:"primarykey" json:"id"`
	Name            string         `gorm:"column:name" json:"name"`
	Email           string         `gorm:"column:email" json:"email"`
	EmailVerifiedAt *time.Time     `gorm:"column:email_verified_at" json:"email_verified_at,omitempty"`
	Password        string         `gorm:"column:password" json:"password"`
	RememberToken   sql.NullString `gorm:"column:remember_token" json:"remember_token,omitempty"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func db() *gorm.DB {
	db, err := gorm.Open(*databases.Resolve(), &gorm.Config{
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("")
			return time.Now().In(utc)
		},
	})

	if errors.Handler("SQL connection error", err) {
		panic(err)
	}

	return db
}

// func (t *Attributes) GetCreatedAt() time.Time {
// 	helpers.DD(t.CreatedAt)
// 	tp, err := time.Parse("2006-01-02 3:04PM", t.CreatedAt)
// 	if err != nil {
// 		logger.Fatal(err)
// 		panic(err)
// 	}
// 	return tp
// }

// func (t *Attributes) GetUpdatedAt() time.Time {
// 	tp, err := time.Parse("2006-01-02 3:04PM", t.CreatedAt)
// 	if err != nil {
// 		logger.Fatal(err)
// 		panic(err)
// 	}
// 	return tp
// }

// func (t *Attributes) GetEmailVerifiedAt() time.Time {
// 	tp, err := time.Parse("2006-01-02 3:04PM", t.CreatedAt)
// 	if err != nil {
// 		logger.Fatal(err)
// 		panic(err)
// 	}
// 	return tp
// }
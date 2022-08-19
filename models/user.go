package models

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uint    `json:"id" gorm:"primary_key;column:id"`
	Name     string  `json:"name" gorm:"column:user_nama"`
	Username string  `json:"username" gorm:"column:user_login"`
	Email    string  `json:"email" gorm:"column:user_email"`
	Picture  *string `json:"picture" gorm:"column:gambar"`
}

func (user *User) AfterFind(tx *gorm.DB) (err error) {
	if user.Picture != nil {
		*user.Picture = fmt.Sprintf("https://content.portalnesia.com/img/content?image=%s", url.QueryEscape(*user.Picture))
	}
	return
}

func (User) TableName() string {
	return fmt.Sprintf("%s_users", os.Getenv("DB_PREFIX"))
}

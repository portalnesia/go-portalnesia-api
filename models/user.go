package models

import (
	"fmt"
	"net/url"

	"gorm.io/gorm"
)

type UserInternal struct {
	ID        uint64  `json:"id" gorm:"primary_key;column:id"`
	Name      string  `json:"name" gorm:"column:user_nama"`
	Username  string  `json:"username" gorm:"column:user_login"`
	Email     *string `json:"email" gorm:"column:user_email"`
	Picture   *string `json:"picture" gorm:"column:gambar"`
	SessionId *string `json:"session_id" gorm:"-"`
}

func (user *UserInternal) AfterFind(tx *gorm.DB) (err error) {
	if user.Picture != nil {
		*user.Picture = fmt.Sprintf("https://content.portalnesia.com/img/content?image=%s", url.QueryEscape(*user.Picture))
	}
	return
}

func (UserInternal) TableName() string {
	return "users"
}

type User struct {
	ID       uint64  `json:"id" gorm:"primary_key;column:id"`
	Name     string  `json:"name" gorm:"column:user_nama"`
	Username string  `json:"username" gorm:"column:user_login"`
	Email    *string `json:"email" gorm:"column:user_email"`
	Picture  *string `json:"picture" gorm:"column:gambar"`
}

func (user *User) AfterFind(tx *gorm.DB) (err error) {
	if user.Picture != nil {
		*user.Picture = fmt.Sprintf("https://content.portalnesia.com/img/content?image=%s", url.QueryEscape(*user.Picture))
	}
	return
}

func (User) TableName() string {
	return "users"
}

type UserUsername struct {
	ID       uint64 `json:"id" gorm:"primary_key;column:id"`
	Username string `json:"username" gorm:"column:user_login"`
}

func (UserUsername) TableName() string {
	return "users"
}

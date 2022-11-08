package models

import (
	"fmt"

	"portalnesia.com/api/config"
)

type Session struct {
	ID        uint64  `json:"id" gorm:"primary_key;column:id"`
	SessionId *string `json:"-" gorm:"column:sess_id"`
	Timestamp *string `json:"-" gorm:"column:timestamp"`
	Pkey      *string `json:"-" gorm:"column:pkey"`
}

func (Session) TableName() string {
	return fmt.Sprintf("%ssession", config.Prefix)
}
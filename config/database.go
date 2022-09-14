package config

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	DBDebug *gorm.DB
	DBProd  *gorm.DB
)

func ChangeDatabase(debug bool) {
	if debug && NODE_ENV != "production" {
		DB = DBDebug
	} else {
		DB = DBProd
	}
}

package models

import (
	"fmt"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"portalnesia.com/api/database"
)

func SetupDB() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	PORT := os.Getenv("DB_PORT")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	var err error
	database.DB, err = gorm.Open("mysql", URL)

	if err != nil {
		panic(err.Error())
	}

	database.DB.AutoMigrate(&User{})
	database.DB.AutoMigrate(&News{})
}

type Timestamp struct {
	Format    string `json:"format" gorm:"-"`
	Timestamp int64  `json:"timestamp" gorm:"-"`
}

package models

import (
	"fmt"
	"log"
	"time"

	"os"

	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"portalnesia.com/api/config"
)

/*func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&News{})
	db.AutoMigrate(&NewsPagination{})
	db.AutoMigrate(&UserContext{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&AccessToken{})
	db.AutoMigrate(&Client{})
}*/

func SetupDB() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASS")
	PORT := os.Getenv("DB_PORT")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")
	level := logger.Error
	if config.NODE_ENV == "test" {
		level = logger.Silent
	} else if config.NODE_ENV == "development" {
		level = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	var err error
	config.DBProd, err = gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: os.Getenv("DB_PREFIX"),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	//defer config.DB.Close()

	//migrateDB(config.DBProd)
}

func SetupDebugDB() {
	USER := os.Getenv("DEBUG_DB_USER")
	PASS := os.Getenv("DB_PASS")
	PORT := os.Getenv("DB_PORT")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DEBUG_DB_NAME")
	level := logger.Error
	if config.NODE_ENV == "test" {
		level = logger.Silent
	}
	//else if config.NODE_ENV == "development" {
	//	level = logger.Info
	// }

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	var err error
	config.DBDebug, err = gorm.Open(mysql.Open(URL), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: os.Getenv("DB_PREFIX"),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	//defer config.DBDebug.Close()

	//migrateDB(config.DBDebug)
}

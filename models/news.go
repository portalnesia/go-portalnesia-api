package models

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/portalnesia/go-utils"
	util "portalnesia.com/api/utils"
)

type News struct {
	Datetime string `json:"-" gorm:"column:datetime"`

	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image" gorm:"column:foto"`
	SourceUrl string    `json:"source_link" gorm:"column:url"`
	Link      string    `json:"link" gorm:"-"`
	Timestamp Timestamp `json:"created" gorm:"-"`
}

type NewsPagination struct {
	Datetime string `json:"-" gorm:"column:datetime"`

	ID        uint      `json:"id" gorm:"primaryKey;column:id"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Image     string    `json:"image" gorm:"column:foto"`
	SourceUrl string    `json:"source_link" gorm:"column:url"`
	Link      string    `json:"link" gorm:"-"`
	Timestamp Timestamp `json:"created" gorm:"-"`
}

func (news *NewsPagination) AfterFind(tx *gorm.DB) (err error) {
	news.Text = util.NewsEncode(news.Text)
	news.Text = utils.CleanAndTruncate(news.Text, 200)

	news.Link = fmt.Sprintf("https://portalnesia.com/news/%s/%s", news.Source, url.QueryEscape(news.Title))
	news.Image = fmt.Sprintf("https://content.portalnesia.com/img/url?image=%s", url.QueryEscape(news.Image))

	date, _ := time.Parse(time.RFC3339, news.Datetime)

	news.Timestamp = Timestamp{
		Timestamp: date.Unix(),
		Format:    utils.TimeAgo(date.Unix()),
	}

	return
}

func (news *News) AfterFind(tx *gorm.DB) (err error) {
	news.Text = util.NewsEncode(news.Text)

	news.Link = fmt.Sprintf("https://portalnesia.com/news/%s/%s", news.Source, url.QueryEscape(news.Title))
	news.Image = fmt.Sprintf("https://content.portalnesia.com/img/url?image=%s", url.QueryEscape(news.Image))

	date, _ := time.Parse(time.RFC3339, news.Datetime)

	news.Timestamp = Timestamp{
		Timestamp: date.Unix(),
		Format:    utils.TimeAgo(date.Unix()),
	}

	return
}

func (News) TableName() string {
	return fmt.Sprintf("%s_news", os.Getenv("DB_PREFIX"))
}

func (NewsPagination) TableName() string {
	return fmt.Sprintf("%s_news", os.Getenv("DB_PREFIX"))
}

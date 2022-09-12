package models

import (
	"fmt"
	"net/url"
	"time"

	"github.com/portalnesia/go-utils"
	"github.com/portalnesia/go-utils/goment"
	"gorm.io/gorm"
	"portalnesia.com/api/config"
	util "portalnesia.com/api/utils"
)

type News struct {
	Datetime time.Time `json:"-" gorm:"column:datetime"`

	ID        uint64               `json:"id" gorm:"primaryKey;column:id"`
	Source    string               `json:"source"`
	Title     string               `json:"title"`
	Text      string               `json:"text"`
	Image     string               `json:"image" gorm:"column:foto"`
	SourceUrl string               `json:"source_link" gorm:"column:url"`
	Link      string               `json:"link" gorm:"-"`
	Timestamp goment.TimeAgoResult `json:"created" gorm:"-"`
}

type NewsPagination struct {
	Datetime time.Time `json:"-" gorm:"column:datetime"`

	ID        uint                 `json:"id" gorm:"primaryKey;column:id"`
	Source    string               `json:"source"`
	Title     string               `json:"title"`
	Text      string               `json:"text"`
	Image     string               `json:"image" gorm:"column:foto"`
	SourceUrl string               `json:"source_link" gorm:"column:url"`
	Link      string               `json:"link" gorm:"-"`
	Timestamp goment.TimeAgoResult `json:"created" gorm:"-"`
}

func (news *NewsPagination) AfterFind(tx *gorm.DB) (err error) {
	news.Text = util.NewsEncode(news.Text)
	news.Text = utils.Clean(news.Text)
	l := len(news.Text)
	ls := 200
	if l < 200 {
		ls = l - 5
	}
	news.Text = utils.Truncate(news.Text, ls)

	news.Link = fmt.Sprintf("https://portalnesia.com/news/%s/%s", news.Source, url.QueryEscape(news.Title))
	news.Image = fmt.Sprintf("https://content.portalnesia.com/img/url?image=%s", url.QueryEscape(news.Image))

	date, _ := utils.NewGoment(news.Datetime)
	news.Timestamp = date.TimeAgo()

	return
}

func (news *News) AfterFind(tx *gorm.DB) (err error) {
	news.Text = util.NewsEncode(news.Text)

	news.Link = fmt.Sprintf("https://portalnesia.com/news/%s/%s", news.Source, url.QueryEscape(news.Title))
	news.Image = fmt.Sprintf("https://content.portalnesia.com/img/url?image=%s", url.QueryEscape(news.Image))

	date, _ := utils.NewGoment(news.Datetime)
	news.Timestamp = date.TimeAgo()

	return
}

func (News) TableName() string {
	return fmt.Sprintf("%snews", config.Prefix)
}

func (NewsPagination) TableName() string {
	return fmt.Sprintf("%snews", config.Prefix)
}

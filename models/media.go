package models

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/portalnesia/go-utils"
	"github.com/portalnesia/go-utils/goment"
	"gorm.io/gorm"
	"portalnesia.com/api/config"
	"portalnesia.com/api/response"
	util "portalnesia.com/api/utils"
)

type Media struct {
	IDNumber      *uint64                `json:"id_number" gorm:"column:id;primaryKey"`
	ID            string                 `json:"id" gorm:"column:unik"`
	Thumbnail     *string                `json:"thumbnail" gorm:"-"`
	Path          *string                `json:"-"`
	Sumber        *string                `json:"-"`
	Thumbs        *string                `json:"-"`
	Private       bool                   `json:"private" gorm:"-"`
	UserID        uint64                 `json:"-" gorm:"column:userid"`
	Downloadtoken *string                `json:"download_token" gorm:"-"`
	Block         bool                   `json:"block" gorm:"column:block"`
	Dilihat       int64                  `json:"-"`
	Title         string                 `json:"title" gorm:"column:judul"`
	Jenis         string                 `json:"-"`
	Type          string                 `json:"type" gorm:"-"`
	Tanggal       time.Time              `json:"-"`
	Artist        *string                `json:"artist"`
	Size          uint64                 `json:"size"`
	URL           string                 `json:"url" gorm:"-"`
	Seen          utils.NumberFormatType `json:"seen" gorm:"-"`
	Created       goment.TimeAgoResult   `json:"created" gorm:"-"`
	User          *User                  `json:"user"`
}

type MyMedia struct {
	Media
	CanSetProfile    bool `json:"can_set_profile" gorm:"-"`
	IsProfilePicture bool `json:"is_profile_picture" gorm:"-"`
}

func (MyMedia) TableName() string {
	return fmt.Sprintf("%sfile", config.Prefix)
}

func (Media) TableName() string {
	return fmt.Sprintf("%sfile", config.Prefix)
}

func (m *Media) JenisToType() {
	if m.Jenis == "lagu" {
		m.Type = "audio"
	} else if m.Jenis == "vdeo" {
		if m.Sumber != nil && *m.Sumber == "youtube" {
			m.Type = "youtube"
		} else {
			m.Type = "video"
		}
	} else {
		m.Type = "images"
	}

	if m.Jenis != "lagu" {
		m.Artist = nil
	}
}

func (m *Media) extractThumbs() (err error) {
	var thumb string
	if m.Type == "images" && m.Path != nil {
		if m.Sumber != nil && *m.Sumber == "imgur" {
			thumb = util.StaticUrl(fmt.Sprintf("img/url?image=%s", url.QueryEscape(*m.Path)))
		} else {
			thumb = util.StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape(*m.Path)))
		}
	} else if m.Type == "audio" {
		if m.Thumbs == nil {
			thumb = util.StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape("images/lagu.png")))
		} else {
			thumb = util.StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape(*m.Thumbs)))
		}
	} else {
		if m.Thumbs != nil {
			thumb = util.StaticUrl(fmt.Sprintf("img/url?image=%s", url.QueryEscape(*m.Thumbs)))
		}
	}
	if thumb != "" {
		m.Thumbnail = &thumb
	}
	return
}

func stripSlashes(s string) string {
	return regexp.MustCompile(`\\`).ReplaceAllString(s, "")
}

func (m *Media) extractUrl() {
	var urls string
	if m.Jenis == "lagu" || m.Jenis == "vdeo" {
		if m.Path != nil && (m.Sumber != nil && *m.Sumber == "firebase") || m.Jenis == "vdeo" {
			urls = stripSlashes(*m.Path)
		} else {
			urls = util.StaticUrl(url.QueryEscape("images/04 Fix You.mp3"))
		}
	} else {
		if m.Sumber != nil && *m.Sumber == "imgur" {
			urls = util.StaticUrl(fmt.Sprintf("img/url?image=%s", url.QueryEscape(*m.Path)))
		} else {
			urls = util.StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape(*m.Path)))
		}
	}
	m.URL = urls
}

func (m *Media) AfterFind(g *gorm.DB) (err error) {
	m.JenisToType()
	m.extractThumbs()
	m.extractUrl()

	date, _ := goment.New(m.Tanggal)
	m.Created = date.TimeAgo()
	if m.Jenis == "lagu" {
		if m.Artist == nil {
			m.Artist = &m.Title
		}
	}
	m.Seen = utils.NumberFormatShort(m.Dilihat)

	return
}

func (m *Media) FormatPagination(c *Context) (err error) {
	var tipe string
	if m.Jenis == "lagu" {
		if m.Sumber != nil && *m.Sumber == "firebase" {
			tipe = "firebase"
		} else {
			tipe = "lagu"
		}
	} else if m.Jenis == "foto" {
		if m.User.ID == c.User.ID {
			tipe = "foto"
		} else {
			tipe = "twibbon"
		}
	}
	if tipe != "" {
		tkn := util.DownloadToken(m.ID, tipe)
		m.Downloadtoken = &tkn
	}
	return
}

func (m *MyMedia) FormatPagination(c *Context) {
	m.Media.FormatPagination(c)

	m.CanSetProfile = m.Sumber == nil
	m.IsProfilePicture = m.Path != nil && c.User.Gambar != nil && *m.Path == *c.User.Gambar
}

func (m *Media) Format(c *Context) (err error) {
	if m.Block {
		if c.User != nil && c.User.ID == m.User.ID {
			return response.Block("files", m.ID, "")
		} else {
			return response.NotFound("files", m.ID, "")
		}
	}
	return
}
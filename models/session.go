package models

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Session struct {
	ID        uint64  `json:"id" gorm:"primary_key;column:id"`
	SessionId *string `json:"-" gorm:"column:sess_id"`
	Timestamp *string `json:"-" gorm:"column:timestamp"`
	Pkey      *string `json:"-" gorm:"column:pkey"`
}

func (Session) TableName() string {
	return "session"
}

type Client struct {
	ClientId      string    `json:"client_id" gorm:"primary_key;column:client_id"`
	ScopeDatabase string    `json:"-" gorm:"column:scope"`
	Scope         []string  `json:"scope" gorm:"-"`
	Grants        string    `json:"grant_types" gorm:"column:grant_types"`
	Internal      bool      `json:"-" gorm:"column:internal"`
	Publish       bool      `json:"publish" gorm:"column:publish"`
	TestUserIdDb  *string   `json:"-" gorm:"column:test_user_id"`
	Block         bool      `json:"block" gorm:"column:block"`
	Origin        *string   `json:"origin" gorm:"column:origin"`
	RedirectUri   *string   `json:"redirect_uri" gorm:"column:redirect_uri"`
	TestUserId    *[]uint64 `json:"test_user_id" gorm:"-"`
	UserId        uint64    `json:"userid" gorm:"column:user_id"`
}

func (client *Client) AfterFind(_ *gorm.DB) (err error) {
	client.Scope = strings.Split(client.ScopeDatabase, " ")
	if client.TestUserIdDb != nil {
		var id []uint64
		for _, ids := range strings.Split(*client.TestUserIdDb, " ") {
			i, _ := strconv.Atoi(ids)
			id = append(id, uint64(i))
		}
		client.TestUserId = &id
	}
	return
}
func (Client) TableName() string {
	return "oauth_clients"
}

type UserContext struct {
	ID              uint64  `json:"id" gorm:"column:id"`
	Name            string  `json:"name" gorm:"column:user_nama"`
	Username        string  `json:"username" gorm:"column:user_login"`
	Email           string  `json:"email" gorm:"column:user_email"`
	Picture         *string `json:"picture" gorm:"column:gambar"`
	SessionId       *string `json:"-" gorm:"column:sess_id"`
	SessionIdNumber *uint   `json:"-" gorm:"column:session_id_number"`
	Timestamp       *string `json:"-" gorm:"column:session_timestamp"`
	Pkey            *string `json:"-" gorm:"column:security_key"`
}

func (UserContext) TableName() string {
	return "users"
}

type AccessToken struct {
	AccessToken   string    `json:"-" gorm:"primary_key;column:access_token"`
	ClientId      *string   `json:"client_id" gorm:"column:client_id"`
	UserId        *uint64   `json:"user_id" gorm:"column:user_id"`
	Expires       *string   `json:"expires" gorm:"column:expires"`
	GrantTypes    *string   `json:"grant_types" gorm:"column:grant_type"`
	ScopeDatabase *string   `json:"-" gorm:"column:scope"`
	Scope         *[]string `json:"scope" gorm:"-"`
}

func (AccessToken) TableName() string {
	return "oauth_access_tokens"
}
func (token *AccessToken) AfterFind(_ *gorm.DB) (err error) {
	var gr []string
	if token.ScopeDatabase != nil {
		gr = strings.Split(*token.ScopeDatabase, " ")
		token.Scope = &gr
	}
	return
}

type ClientContext struct {
	ClientId    string
	Scope       *[]string
	Grants      string
	AccessToken *string
}

type Context struct {
	// IS DEVELOPER CLIENT (EXTERNAL APPLICATION WITH OAUTH2)
	IsApi bool
	// IS WEB APPLICATION
	IsWeb bool
	// IS INTERNAL PORTALNESIA
	IsInternal bool
	// IS ACCESSED FROM LOCALHOST PORTALNESIA
	IsDebug bool
	// IS USE NATIVE APPLICATION
	IsNative bool
	// IS ACCESSED FROM PHP OR INTERNAL SERVER
	IsInternalServer bool
	// CHECKLIST FOR FIRST AUTHORIZATION AND SECOND
	Checklist     bool
	AlmostExpired bool
	User          *UserContext
	// CLIENT OBJECT IF WITH OAUTH2 ACCESS TOKEN
	Client *ClientContext
}

var CtxDefaultValue = Context{
	IsApi:            false,
	IsWeb:            false,
	IsInternal:       false,
	IsDebug:          false,
	IsNative:         false,
	IsInternalServer: false,
	Checklist:        false,
	AlmostExpired:    false,
}

type ContextUserConfig struct {
	WithEmail bool
	SessionId *string
}

func (c *Context) ToUserModels(g *gorm.DB, config ContextUserConfig) *User {
	if c == nil {
		return nil
	} else if c.User == nil {
		return nil
	} else {
		user := &User{
			ID:       c.User.ID,
			Name:     c.User.Name,
			Username: c.User.Username,
			Picture:  c.User.Picture,
		}
		if config.WithEmail {
			user.Email = &c.User.Email
		}

		user.AfterFind(g)
		return user
	}
}
func (c *Context) ToUserInternalModels(g *gorm.DB, config ContextUserConfig) *UserInternal {
	if c == nil {
		return nil
	} else if c.User == nil {
		return nil
	} else {
		user := &UserInternal{
			ID:       c.User.ID,
			Name:     c.User.Name,
			Username: c.User.Username,
			Picture:  c.User.Picture,
		}
		if config.WithEmail {
			user.Email = &c.User.Email
		}
		if config.SessionId != nil {
			user.SessionId = config.SessionId
		}

		user.AfterFind(g)
		return user
	}
}

//start = page > 1 ? (page*per_page)-per_page : 0;

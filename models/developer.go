package models

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"portalnesia.com/api/config"
)

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
	return fmt.Sprintf("%soauth_clients", config.Prefix)
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
	return fmt.Sprintf("%soauth_access_tokens", config.Prefix)
}
func (token *AccessToken) AfterFind(_ *gorm.DB) (err error) {
	var gr []string
	if token.ScopeDatabase != nil {
		gr = strings.Split(*token.ScopeDatabase, " ")
		token.Scope = &gr
	}
	return
}
package models

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"gorm.io/gorm"
	"portalnesia.com/api/config"
)

type WebauthnKeys struct {
	// A probabilistically-unique byte sequence identifying a public key credential source and its authentication assertions.
	ID []byte `json:"id"`
	// The public key portion of a Relying Party-specific credential key pair, generated by an authenticator and returned to
	// a Relying Party at registration time (see also public key credential). The private key portion of the credential key
	// pair is known as the credential private key. Note that in the case of self attestation, the credential key pair is also
	// used as the attestation key pair, see self attestation for details.
	PublicKey string    `json:"key"`
	Datetime  time.Time `json:"datetime"`
	Device    string    `json:"device"`
}

type UserWebauthn struct {
	// API KEY
	ID           string         `json:"id"`
	UserID       uint64         `json:"user_id"`
	WebauthnKeys []WebauthnKeys `json:"webauthnkeys"`
}

type User struct {
	ID           uint64       `json:"id" gorm:"primaryKey;column:id"`
	Name         string       `json:"name" gorm:"column:user_nama"`
	Username     string       `json:"username" gorm:"column:user_login"`
	Email        *string      `json:"email" gorm:"column:user_email"`
	Gambar       *string      `json:"-" gorm:"column:gambar"`
	Picture      *string      `json:"picture" gorm:"-"`
	Private      bool         `json:"private" gorm:"column:private"`
	MediaPrivate bool         `json:"media_private" gorm:"column:media_private"`
	Verify       bool         `json:"verify" gorm:"column:verify"`
	Paid         bool         `json:"-" gorm:"column:paid"`
	PaidExpired  string       `json:"-" gorm:"column:paid_expired"`
	Webauthn     UserWebauthn `json:"-" gorm:"column:security_key"`
}

func (user *User) AfterFind(tx *gorm.DB) (err error) {
	if user.Gambar != nil {
		pic := fmt.Sprintf("https://content.portalnesia.com/img/content?image=%s", url.QueryEscape(*user.Gambar))
		user.Picture = &pic
	}
	return
}

func (User) TableName() string {
	return fmt.Sprintf("%susers", config.Prefix)
}

func (User) TableFollowName() string {
	return fmt.Sprintf("%sfollow", config.Prefix)
}

type UserInternal struct {
	User
	SessionId *string `json:"session_id" gorm:"-"`
}

func (user *UserInternal) AfterFind(tx *gorm.DB) (err error) {
	user.User.AfterFind(tx)

	return
}

func (UserInternal) TableName() string {
	return fmt.Sprintf("%susers", config.Prefix)
}

type UserDetail struct {
	User
	IsFollowing     bool `json:"isFollowing" gorm:"-"`
	IsFollower      bool `json:"isFollower" gorm:"-"`
	IsFollowPending bool `json:"isFollowPending" gorm:"-"`
}

func (user *UserDetail) AfterFind(g *gorm.DB) (err error) {
	user.User.AfterFind(g)
	return
}

func (UserDetail) TableName() string {
	return fmt.Sprintf("%susers", config.Prefix)
}

type UserUsername struct {
	ID       uint64 `json:"id" gorm:"primary_key;column:id"`
	Username string `json:"username" gorm:"column:user_login"`
}

func (UserUsername) TableName() string {
	return fmt.Sprintf("%susers", config.Prefix)
}

func (u *UserDetail) Format(c *UserContext) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	pending, follower, following := make(chan bool), make(chan bool), make(chan bool)

	go func() {
		defer wg.Done()
		defer close(pending)

		pending <- u.CheckIsFollowPendings(c)
	}()
	go func() {
		defer wg.Done()
		defer close(follower)

		follower <- u.CheckIsFollowers(c)
	}()
	go func() {
		defer wg.Done()
		defer close(following)

		following <- u.CheckIsFollowings(c)
	}()

	u.IsFollowPending = <-pending
	u.IsFollower = <-follower
	u.IsFollowing = <-following
	wg.Wait()
}

func (u *UserDetail) CheckIsFollowings(user *UserContext) bool {
	var exists bool
	db := config.DB
	db.Table(u.TableFollowName()).Select("count(*) > 0").Where("difollow = ? AND yangfollow = ?", u.ID, user.ID).Session(&gorm.Session{}).Find(&exists)

	return exists
}
func (u *UserDetail) CheckIsFollowers(user *UserContext) bool {
	exists := false
	db := config.DB
	db.Table(u.TableFollowName()).Select("count(*) > 0").Where("difollow = ? AND yangfollow = ? AND pending='0'", user.ID, u.ID).Session(&gorm.Session{}).Find(&exists)
	return exists
}
func (u *UserDetail) CheckIsFollowPendings(user *UserContext) bool {
	exists := false
	db := config.DB
	db.Table(u.TableFollowName()).Select("count(*) > 0").Where("difollow = ? AND yangfollow = ? AND pending='1'", u.ID, user.ID).Session(&gorm.Session{}).Find(&exists)
	return exists
}

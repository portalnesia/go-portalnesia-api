package utils

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"portalnesia.com/api/config"
)

func TestStaticUrl(t *testing.T) {
	parse := StaticUrl("")
	parse2 := StaticUrl("img/content?image=Banner.png")

	if parse != os.Getenv("STATIC_URL") {
		t.Errorf("StaticUrl: ``, Get %s", parse)
	}
	if parse2 != fmt.Sprintf("%s/img/content?image=Banner.png", os.Getenv("STATIC_URL")) {
		t.Errorf("StaticUrl: `img/content?image=Banner.png`, Get %s", parse2)
	}
}

func TestHref(t *testing.T) {
	parse := Href("")
	parse2 := Href("v1/user")

	if parse != os.Getenv("API_URL") {
		t.Errorf("Href: ``, Get %s", parse)
	}
	if parse2 != fmt.Sprintf("%s/v1/user", os.Getenv("API_URL")) {
		t.Errorf("Href: `v1/user`, Get %s", parse2)
	}
}

func TestLinkUrl(t *testing.T) {
	parse := LinkUrl("")
	parse2 := LinkUrl("v1/user")

	if parse != os.Getenv("LINK_URL") {
		t.Errorf("LinkUrl: ``, Get %s", parse)
	}
	if parse2 != fmt.Sprintf("%s/v1/user", os.Getenv("LINK_URL")) {
		t.Errorf("LinkUrl: `v1/user`, Get %s", parse2)
	}
}

func TestAccountUrl(t *testing.T) {
	parse := AccountUrl("")
	parse2 := AccountUrl("login")

	if parse != os.Getenv("ACCOUNT_URL") {
		t.Errorf("AccountUrl: ``, Get %s", parse)
	}
	if parse2 != fmt.Sprintf("%s/login", os.Getenv("ACCOUNT_URL")) {
		t.Errorf("AccountUrl: `login`, Get %s", parse2)
	}
}

func TestPortalUrl(t *testing.T) {
	parse := PortalUrl("")
	parse2 := PortalUrl("contact")

	if parse != os.Getenv("PORTAL_URL") {
		t.Errorf("PortalUrl: ``, Get %s", parse)
	}
	if parse2 != fmt.Sprintf("%s/contact", os.Getenv("PORTAL_URL")) {
		t.Errorf("PortalUrl: `contact`, Get %s", parse2)
	}
}

func TestAnalyzeStaticUrl(t *testing.T) {
	path := "Banner.png"
	result := AnalyzeStaticUrl(path)

	if result != fmt.Sprintf("%s/img/content?image=%s", os.Getenv("STATIC_URL"), url.QueryEscape(path)) {
		t.Errorf("AnalyzeStaticUrl: `Not URL`, Get %s", result)
	}

	path = "https://picsum.photos/200"
	result = AnalyzeStaticUrl(path)
	if result != fmt.Sprintf("%s/img/url?image=%s", os.Getenv("STATIC_URL"), url.QueryEscape(path)) {
		t.Errorf("AnalyzeStaticUrl: `URL`, Get %s", result)
	}
}

func TestProfileUrl(t *testing.T) {
	path := "Banner.png"
	result := ProfileUrl(&path)

	if *result != fmt.Sprintf("%s/img/content?image=%s", os.Getenv("STATIC_URL"), url.QueryEscape(path)) {
		t.Errorf("AnalyzeStaticUrl: `URL`, Get %s", *result)
	}

	result = ProfileUrl(nil)
	if result != nil {
		t.Errorf("AnalyzeStaticUrl: `Nil`, Get %s", *result)
	}
}

func getToken(key string, secret string, date string) TokenBase {
	return TokenBase{
		Token:    secret,
		Key:      key,
		Date:     date,
		Datetime: date,
	}
}

func TestVerifyToken(t *testing.T) {
	config.SetupConfig()

	token_secret := "This is secret"
	token_date := time.Now().Format(time.RFC3339)
	token_date_false := time.Now().Add(-time.Minute * 61).Format(time.RFC3339)
	// Valid All
	encrypted := CreateToken(getToken("", token_secret, token_date))
	verify := VerifyToken[TokenBase](encrypted, token_secret, int64(time.Hour)*1)

	if verify.Data.Date != token_date {
		t.Errorf("[VerifyToken1] Invalid token date, Get %s", verify.Data.Date)
	}
	if verify.Data.Token != token_secret {
		t.Errorf("[VerifyToken1] Invalid token secret, Get %s", verify.Data.Token)
	}
	if !verify.Info.Date {
		t.Errorf("[VerifyToken1] Token date is not verified")
	}
	if !verify.Info.Token {
		t.Errorf("[VerifyToken1] Token key is not verified")
	}
	if !verify.Verified {
		t.Errorf("[VerifyToken1] Token is not verified, Get %s", verify.Data.Token)
	}

	// Invalid Date
	encrypted = CreateToken(getToken("", token_secret, token_date_false))
	verify = VerifyToken[TokenBase](encrypted, token_secret, int64(time.Hour)*1)
	if verify.Data.Date != token_date_false {
		t.Errorf("[VerifyToken2] Invalid token date, Get %s", verify.Data.Date)
	}
	if verify.Data.Token != token_secret {
		t.Errorf("[VerifyToken2] Invalid token secret, Get %s", verify.Data.Token)
	}
	if !verify.Info.Token {
		t.Errorf("[VerifyToken2] Token key is not verified")
	}
	if verify.Info.Date {
		t.Errorf("[VerifyToken2] Token date is verified (false)")
	}
	if verify.Verified {
		t.Errorf("[VerifyToken2] Token is verified (false)")
	}

	// Invalid Key
	key := "invalid key"
	encrypted = CreateToken(getToken("", key, token_date))
	verify = VerifyToken[TokenBase](encrypted, token_secret, int64(time.Hour)*1)
	if verify.Data.Date != token_date {
		t.Errorf("[VerifyToken3] Invalid token date, Get %s", verify.Data.Date)
	}
	if verify.Data.Token != key {
		t.Errorf("[VerifyToken3] Invalid token secret, Get %s", verify.Data.Token)
	}
	if verify.Info.Date {
		t.Errorf("[VerifyToken3] Token date is verified (false)")
	}
	if verify.Verified {
		t.Errorf("[VerifyToken3] Token is verified (false)")
	}
}

func TestVerifyTokenAuth(t *testing.T) {
	config.SetupConfig()

	token_date := time.Now().Format(time.RFC3339)
	token_date_false := time.Now().AddDate(0, 0, -32).Format(time.RFC3339)

	key := fmt.Sprintf("$2a$08$j9jNyZvS.KFPHIMRAEE4k.ckWmeTMdv17E3QvftgbxEfAO0K94nDm%s", os.Getenv("DEBUG_USERID"))

	// Valid All
	encrypted := CreateToken(getToken(key, "", token_date))
	verify := VerifyTokenAuth(encrypted)

	if verify.Data.Date != token_date {
		t.Errorf("[VerifyToken1] Invalid token date, Get %s", verify.Data.Date)
	}
	if verify.Data.Key != key {
		t.Errorf("[VerifyToken1] Invalid token secret, Get %s", verify.Data.Token)
	}
	if !verify.Info.Date {
		t.Errorf("[VerifyToken1] Token date is not verified")
	}
	if !verify.Verified {
		t.Errorf("[VerifyToken1] Token is not verified, Get %s", verify.Data.Token)
	}

	// Invalid Date
	encrypted = CreateToken(getToken(key, "", token_date_false))
	verify = VerifyTokenAuth(encrypted)
	if verify.Data.Date != token_date_false {
		t.Errorf("[VerifyToken2] Invalid token date, Get %s", verify.Data.Date)
	}
	if verify.Data.Key != key {
		t.Errorf("[VerifyToken2] Invalid token secret, Get %s", verify.Data.Token)
	}
	if verify.Info.Date {
		t.Errorf("[VerifyToken2] Token date is verified (false)")
	}
	if verify.Verified {
		t.Errorf("[VerifyToken2] Token is verified (false)")
	}
}

func TestItemExists(t *testing.T) {
	data_string := []string{"hello", "world", "this", "is", "from", "testing"}
	data_int := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	ok := ItemExists(data_string, "this")
	if !ok {
		t.Errorf("[ItemExists] string, Should be exists.\nArray: %+v\nItem: %s", data_string, "this")
	}

	ok = ItemExists(data_string, "notfound")
	if ok {
		t.Errorf("[ItemExists] string, Should be not exists.\nArray: %+v\nItem: %s", data_string, "notfound")
	}

	ok = ItemExists(data_int, 7)
	if !ok {
		t.Errorf("[ItemExists] int, Should be exists.\nArray: %+v\nItem: %s", data_int, "this")
	}

	ok = ItemExists(data_int, 12)
	if ok {
		t.Errorf("[ItemExists] int, Should be not exists.\nArray: %+v\nItem: %s", data_int, "notfound")
	}
}

func TestCheckGrants(t *testing.T) {
	grant := []string{"authorization_code"}
	db_grant := []string{"authorization_code", "refresh_tokens", "client_credentials"}

	ok := CheckGrants(grant, db_grant)
	if !ok {
		t.Errorf("[CheckGrants] Should be true.\nGrants: %+v\nDB Grants: %s", grant, db_grant)
	}

	db_grant = []string{"refresh_tokens", "client_credentials"}
	ok = CheckGrants(grant, db_grant)
	if ok {
		t.Errorf("[CheckGrants] Should be false.\nGrants: %+v\nDB Grants: %s", grant, db_grant)
	}
}

func TestCheckScope(t *testing.T) {
	client_scope := []string{"email", "profile", "openid"}
	checked_scope := []string{"profile"}

	ok := CheckScope(client_scope, checked_scope)
	if !ok {
		t.Errorf("[CheckScope] Should be true.\nClient Scope: %+v\nChecked Scope: %s", client_scope, checked_scope)
	}

	checked_scope = []string{"news"}
	ok = CheckScope(client_scope, checked_scope)
	if ok {
		t.Errorf("[CheckScope] Should be false.\nClient Scope: %+v\nChecked Scope: %s", client_scope, checked_scope)
	}

	client_scope = []string{"superuser"}
	ok = CheckScope(client_scope, checked_scope)
	if !ok {
		t.Errorf("[CheckScope] Should be true.\nClient Scope: %+v\nChecked Scope: %s", client_scope, checked_scope)
	}
}

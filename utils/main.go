package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gofiber/fiber/v2"
	"github.com/portalnesia/go-utils"
	"github.com/portalnesia/go-utils/goment"
	"portalnesia.com/api/config"
)

func parsePath(path string) string {
	if path != "" {
		path = fmt.Sprintf("/%s", path)
	}
	return path
}

func StaticUrl(path string) string {
	return fmt.Sprintf("%s%s", os.Getenv("STATIC_URL"), parsePath(path))
}

func Href(path string) string {
	return fmt.Sprintf("%s%s", os.Getenv("API_URL"), parsePath(path))
}

func LinkUrl(path string) string {
	return fmt.Sprintf("%s%s", os.Getenv("LINK_URL"), parsePath(path))
}

func AccountUrl(path string) string {
	return fmt.Sprintf("%s%s", os.Getenv("ACCOUNT_URL"), parsePath(path))
}

func PortalUrl(path string) string {
	return fmt.Sprintf("%s%s", os.Getenv("PORTAL_URL"), parsePath(path))
}

func AnalyzeStaticUrl(path string) string {
	if utils.IsUrl(path) {
		return StaticUrl(fmt.Sprintf("img/url?image=%s", url.QueryEscape(path)))
	} else {
		return StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape(path)))
	}
}

func ProfileUrl(path *string) *string {
	p := "images/avatar.png"
	if path == &p || path == nil {
		return nil
	} else {
		p = AnalyzeStaticUrl(*path)
		return &p
	}
}

type TokenBase struct {
	Token    string `json:"token"`
	Date     string `json:"date"`
	Datetime string `json:"datetime"`
	Key      string `json:"key"`
}
type TokenVerifiedInfo struct {
	Token bool
	Date  bool
}
type TokenResponse[T any] struct {
	Verified bool
	Data     T
	Date     time.Time
	Info     TokenVerifiedInfo
}

func VerifyToken[T any](datatoken string, secret string, second int64) TokenResponse[T] {
	decryptString, err := config.Crypto.Decrypt(datatoken)
	var res T

	result := TokenResponse[T]{
		Verified: false,
		Data:     res,
		Date:     time.Time{},
		Info: TokenVerifiedInfo{
			Token: false,
			Date:  false,
		},
	}

	if err != nil || decryptString == "" {
		return result
	}

	err = json.Unmarshal([]byte(decryptString), &res)
	if err != nil {
		return result
	}
	result.Data = res
	st := reflect.ValueOf(res)
	f := st.FieldByName("Token")
	if f.IsZero() {
		return result
	}
	t := f.Interface().(string)

	if t != secret {
		return result
	}
	result.Info.Token = true

	f = st.FieldByName("Date")
	if f.IsZero() {
		return result
	}
	t = f.Interface().(string)

	d, err := dateparse.ParseAny(t)
	if err != nil {
		return result
	}
	result.Date = d
	tn := d.Add(time.Duration(second)).After(time.Now())
	if !tn {
		return result
	}
	result.Info.Date = true
	result.Verified = true
	return result
}

func VerifyTokenAuth(datatoken string) TokenResponse[TokenBase] {
	decryptString, err := config.Crypto.Decrypt(datatoken)

	result := TokenResponse[TokenBase]{
		Verified: false,
		Data:     TokenBase{},
		Date:     time.Time{},
		Info: TokenVerifiedInfo{
			Token: true,
			Date:  false,
		},
	}

	if err != nil || decryptString == "" {
		return result
	}
	var res TokenBase
	err = json.Unmarshal([]byte(decryptString), &res)
	if err != nil {
		return result
	}
	result.Data = res
	st := reflect.ValueOf(res)

	f := st.FieldByName("Datetime")
	if f.IsZero() {
		return result
	}
	t := f.Interface().(string)

	d, err := dateparse.ParseAny(t)
	if err != nil {
		return result
	}
	result.Date = d
	tn := d.AddDate(0, 0, 30).After(time.Now())
	if !tn {
		return result
	}
	result.Info.Date = true
	result.Verified = true
	return result
}

func CreateToken(data interface{}) string {
	dt, _ := json.Marshal(data)
	encrypted, _ := config.Crypto.Encrypt(string(dt))
	return encrypted
}

func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func CheckGrants(grant []string, db_grant []string) bool {
	for _, g := range db_grant {
		if ItemExists(grant, g) {
			return true
		}
	}
	return false
}

func CheckScope(client_scope []string, scope []string) bool {
	scope = append(scope, "superuser")
	for _, s := range client_scope {
		if ItemExists(scope, s) {
			return true
		}
	}
	return false
}

type PartialFiberCookie struct {
	Name    string
	Value   string
	Expires goment.PortalnesiaGoment
}

func SetCookie(c *fiber.Ctx, cookie PartialFiberCookie) {
	u, _ := utils.ParseUrl(os.Getenv("PORTAL_URL"))
	var d string
	if config.IsProduction {
		d = fmt.Sprintf(".%s", u)
	} else {
		d = "localhost"
	}
	c.Cookie(&fiber.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Expires:  cookie.Expires.ToTime(),
		Domain:   d,
		HTTPOnly: false,
		Secure:   config.IsProduction,
	})
}

type DownloadTokenStruct struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	TypeID string `json:"type_id"`
}

func DownloadToken(id string, tipe string) string {
	date, _ := utils.NewGoment()

	return CreateToken(DownloadTokenStruct{
		ID:     id,
		Date:   date.PNformat(),
		TypeID: tipe,
	})
}

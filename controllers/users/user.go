package user_controllers

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/portalnesia/go-utils"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"
	util "portalnesia.com/api/utils"
)

// v1/user
func FindMe(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)

	if ctx.User != nil { // Is logged in ?
		if ctx.IsWeb { // Is web application ?
			sess := ctx.User.SessionId
			if sess == nil {
				sess_id := utils.NanoId()
				sess = &sess_id
				db.Table("session").Where("id = ?", &ctx.User.SessionIdNumber).Update("sess_id", sess_id)
			}
			t, _ := utils.NewGoment()
			t.Add(30, "days")
			if sess != nil {
				util.SetCookie(c, util.PartialFiberCookie{
					Name:    "pn_sess",
					Value:   *sess,
					Expires: *t,
				})
			}

			if ctx.User.Timestamp != nil {
				t, _ = utils.NewGoment()
				tt, err := utils.NewGoment(*ctx.User.Timestamp)
				if err == nil {
					tt.Add(20, "days")
					if t.IsAfter(tt.Goment) {
						portalid := c.Cookies("portalid", "")
						portalid = strings.ReplaceAll(portalid, "%3A", ":")
						db.Table("session").Where("id = ?", &ctx.User.SessionIdNumber).Update("sess_time", t.PNformat())
						auth := util.VerifyTokenAuth(portalid)
						if auth.Data != nil && auth.Data.Key != nil {
							json := util.CreateToken(map[string]interface{}{
								"key":      *auth.Data.Key,
								"datetime": "2025-05-05 20:20:00",
							})
							t.Add(30, "days")
							util.SetCookie(c, util.PartialFiberCookie{
								Name:    "portalid",
								Value:   json,
								Expires: *t,
							})
						}
					}
				}
			}
		}

		config := models.ContextUserConfig{
			WithEmail: ctx != nil && (ctx.IsWeb || ctx.Client != nil && ctx.Client.Scope != nil && util.CheckScope(*ctx.Client.Scope, []string{"email"})),
			SessionId: ctx.User.SessionId,
		}
		var users interface{}
		if ctx.IsInternal {
			users = ctx.ToUserInternalModels(db, config)
		} else {
			users = ctx.ToUserModels(db, config)
		}
		return response.Response(users).Send(c)
	} else { // Not logged in
		if ctx.IsWeb { // Is web application
			t, _ := utils.NewGoment()
			t.Add(30, "days")
			pn_sess := c.Cookies("pn_sess", "")
			if pn_sess != "" {
				sess_id := utils.NanoId()
				util.SetCookie(c, util.PartialFiberCookie{
					Name:    "pn_sess",
					Value:   sess_id,
					Expires: *t,
				})
			} else { // Other application
				var sess models.Session
				if err := db.First(&sess, "sess_id = ?", pn_sess).Error; err == nil {
					if sess.SessionId != nil {
						util.SetCookie(c, util.PartialFiberCookie{
							Name:    "pn_sess",
							Value:   *sess.SessionId,
							Expires: *t,
						})
					}
				}
			}
		}
	}

	return response.Response[*models.User](nil).Send(c)
}

// /v1/user/list
func ListUsername(c *fiber.Ctx) error {
	db := config.DB
	tipe := c.Query("type", "username")
	if tipe != "username" && tipe != "all" {
		tipe = "username"
	}
	q := c.Query("q", "")
	var user models.UserUsername
	if q != "" {
		q, _ = url.QueryUnescape(q)
		db.Where("user_login LIKE ? AND active='1' AND remove='0' AND block='0' AND suspend='0'", q).Group("user_login").Limit(100).Find(&user)
	} else {
		db.Where("active='1' AND remove='0' AND block='0' AND suspend='0'").Group("user_login").Limit(100).Find(&user)
	}
	return response.Response(user).Send(c)
}

// /v1/user/:username
func FindUser(c *fiber.Ctx) error {
	db := config.DB
	id := c.Params("id", "")

	var user models.User

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	return response.Response(user).Send(c)
}

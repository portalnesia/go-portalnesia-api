package user_controllers

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

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
						if auth.Data.Key != "" {
							json := util.CreateToken(map[string]interface{}{
								"key":      auth.Data.Key,
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
	var user []models.UserUsername
	if q != "" {
		q, _ = url.QueryUnescape(q)
		db.Where("user_login LIKE ? AND active='1' AND remove='0' AND block='0' AND suspend='0'", fmt.Sprintf("%%%s%%", q)).Group("user_login").Limit(100).Find(&user)
	} else {
		db.Where("active='1' AND remove='0' AND block='0' AND suspend='0'").Group("user_login").Limit(100).Find(&user)
	}
	return response.Response(user).Send(c)
}

// /v1/user/:id
func FindUser(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)
	id := c.Params("id", "")

	var user models.UserDetail

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}
	if ctx.User.ID != user.ID {
		user.Format(ctx.User)
	}

	return response.Response(user).Send(c)
}

// /v1/user/:id/followers
func FindFollowers(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)
	id := c.Params("id", "")
	var user models.UserDetail

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	if ctx.User == nil && user.Private || ctx.User != nil && ctx.User.ID != user.ID && user.Private && !user.IsFollowing {
		return response.PrivatePagination().Send(c)
	}
	tblUser := user.TableName()
	tblFollow := fmt.Sprintf("%sfollow", config.Prefix)

	g := db.Select(fmt.Sprintf("%s.*", tblUser)).Joins(fmt.Sprintf("JOIN %s on %s.yangfollow = %s.id", tblFollow, tblFollow, tblUser)).Where(fmt.Sprintf("%s.difollow = ? AND %s.pending = '0' AND %s.active='1' AND %s.remove='0' AND %s.block='0' AND %s.suspend='0'", tblFollow, tblFollow, tblUser, tblUser, tblUser, tblUser), user.ID).Order(fmt.Sprintf("%s.tanggal desc", tblFollow))

	resp := response.GetPagination[models.UserDetail](c).PaginationResponse(g, g)
	if len(resp.Data.Data) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(resp.Data.Data))
		for i := range resp.Data.Data {
			go func(i int) {
				defer wg.Done()
				resp.Data.Data[i].Format(ctx.User)
			}(i)
		}
		wg.Wait()
	}
	return resp.Send(c)
}

// /v1/user/:id/followers/pending
func FindFollowersPending(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)
	id := c.Params("id", "")
	var user models.UserDetail

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	if ctx.User == nil || ctx.User.ID != user.ID {
		return response.EndpointNotFound()
	}

	tblUser := user.TableName()
	tblFollow := fmt.Sprintf("%sfollow", config.Prefix)

	g := db.Select(fmt.Sprintf("%s.*", tblUser)).Joins(fmt.Sprintf("JOIN %s on %s.yangfollow = %s.id", tblFollow, tblFollow, tblUser)).Where(fmt.Sprintf("%s.difollow = ? AND %s.pending = '1' AND %s.active='1' AND %s.remove='0' AND %s.block='0' AND %s.suspend='0'", tblFollow, tblFollow, tblUser, tblUser, tblUser, tblUser), ctx.User.ID).Order(fmt.Sprintf("%s.tanggal desc", tblFollow))

	resp := response.GetPagination[models.UserDetail](c).PaginationResponse(g, g)

	if len(resp.Data.Data) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(resp.Data.Data))
		for i := range resp.Data.Data {
			go func(i int) {
				defer wg.Done()
				resp.Data.Data[i].Format(ctx.User)
			}(i)
		}
		wg.Wait()
	}

	return resp.Send(c)
}

// /v1/user/:id/following
func FindFollowings(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)
	id := c.Params("id", "")
	var user models.UserDetail

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	if ctx.User == nil && user.Private || ctx.User != nil && ctx.User.ID != user.ID && user.Private && !user.IsFollowing {
		return response.PrivatePagination().Send(c)
	}
	tblUser := user.TableName()
	tblFollow := fmt.Sprintf("%sfollow", config.Prefix)

	g := db.Select(fmt.Sprintf("%s.*", tblUser)).Joins(fmt.Sprintf("JOIN %s on %s.difollow = %s.id", tblFollow, tblFollow, tblUser)).Where(fmt.Sprintf("%s.yangfollow = ? AND %s.pending = '0' AND %s.active='1' AND %s.remove='0' AND %s.block='0' AND %s.suspend='0'", tblFollow, tblFollow, tblUser, tblUser, tblUser, tblUser), user.ID).Order(fmt.Sprintf("%s.tanggal desc", tblFollow))

	resp := response.GetPagination[models.UserDetail](c).PaginationResponse(g, g)
	if len(resp.Data.Data) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(resp.Data.Data))
		for i := range resp.Data.Data {
			go func(i int) {
				defer wg.Done()
				resp.Data.Data[i].Format(ctx.User)
			}(i)
		}
		wg.Wait()
	}
	return resp.Send(c)
}

// /v1/user/:id/following
func FindMedia(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)
	id := c.Params("id", "")
	var user models.UserDetail

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	user.Format(ctx.User)
	if user.MediaPrivate && (ctx.User == nil || ctx.User.ID != user.ID && !user.IsFollowing) {
		return response.PrivatePagination().Send(c)
	}

	tblFile := fmt.Sprintf("%sfile", config.Prefix)

	if ctx.User.ID == user.ID {
		var test []models.MyMedia

		g := db.Model(&test).Where(fmt.Sprintf("%s.userid = ? AND %s.private = '0' AND %s.jenis != 'apps' AND (%s.jenis = 'lagu' OR %s.jenis = 'foto' OR %s.jenis='vdeo') AND %s.tampil='1' AND %s.block='0'", tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile), user.ID).Order(fmt.Sprintf("%s.tanggal DESC", tblFile)).Preload("User")

		t := db.Table(models.Media{}.TableName()).Where(fmt.Sprintf("%s.userid = ? AND %s.private = '0' AND %s.jenis != 'apps' AND (%s.jenis = 'lagu' OR %s.jenis = 'foto' OR %s.jenis='vdeo') AND %s.tampil='1' AND %s.block='0'", tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile), user.ID)

		resp := response.GetPagination[models.MyMedia](c).PaginationResponse(g, t)
		for i := range resp.Data.Data {
			resp.Data.Data[i].FormatPagination(ctx)
		}
		return resp.Send(c)
	} else {
		var test []models.Media

		g := db.Model(&test).Where(fmt.Sprintf("%s.userid = ? AND %s.private = '0' AND %s.jenis != 'apps' AND (%s.jenis = 'lagu' OR %s.jenis = 'foto' OR %s.jenis='vdeo') AND %s.tampil='1' AND %s.block='0'", tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile), user.ID).Order(fmt.Sprintf("%s.tanggal DESC", tblFile)).Preload("User")

		t := db.Table(models.Media{}.TableName()).Where(fmt.Sprintf("%s.userid = ? AND %s.private = '0' AND %s.jenis != 'apps' AND (%s.jenis = 'lagu' OR %s.jenis = 'foto' OR %s.jenis='vdeo') AND %s.tampil='1' AND %s.block='0'", tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile, tblFile), user.ID)

		resp := response.GetPagination[models.Media](c).PaginationResponse(g, t)
		for i := range resp.Data.Data {
			resp.Data.Data[i].FormatPagination(ctx)
		}
		return resp.Send(c)
	}
}

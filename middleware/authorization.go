package middleware

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"
	util "portalnesia.com/api/utils"
)

type AuthorizationConfig struct {
	Disable bool
}

type AuthInternal struct {
	util.TokenBase
	UserId    *uint64 `json:"userid"`
	SessionId *string `json:"session_id"`
}
type AuthWeb struct {
	util.TokenBase
}

func Authorization(options AuthorizationConfig) func(*fiber.Ctx) error {
	tblSess := fmt.Sprintf("%ssession", config.Prefix)
	tblUser := fmt.Sprintf("%susers", config.Prefix)
	tblAccess := fmt.Sprintf("%soauth_access_tokens", config.Prefix)

	return func(c *fiber.Ctx) error {
		ctx, ok := c.Locals("ctx").(*models.Context)
		xDeviceId := c.Get("x-device-id", "")
		xSessionId := c.Get("x-session-id", "")
		xApplicationVersion := c.Get("x-application-version", "")
		xAppToken := c.Get("x-app-token", "")
		pnAuth := c.Get("pn-auth", "")
		pnInternalPortalnesia := c.Get("pn-internal-portalnesia", "")
		clientId := c.Get("pn-client-id", "")
		xDebug := c.Get("x-debug", "")
		auth := c.Get("authorization", "")
		portalid := c.Cookies("portalid", "")
		portalid = strings.ReplaceAll(portalid, "%3A", ":")
		var client models.Client

		if !ok || ctx == nil {
			ctx = &models.CtxDefaultValue
			if xDeviceId != "" && xSessionId != "" && xApplicationVersion != "" && xAppToken != "" {
				ctx.IsNative = true
				c.Locals("browserStr", fmt.Sprintf("Portalnesia on Android v%s", xApplicationVersion))
				ctx.IsDebug = regexp.MustCompile(`\-debug$`).MatchString(xApplicationVersion)
			} else if pnAuth != "" {
				ctx.IsWeb = true
			} else if pnInternalPortalnesia != "" {
				verify := util.VerifyToken[AuthInternal](pnInternalPortalnesia, os.Getenv("AUTH_PN_INTERNAL_SECRET"), int64(time.Hour)*2)
				if verify.Verified {
					if verify.Data.UserId != nil && verify.Data.SessionId != nil {
						var user models.UserContext
						sel := fmt.Sprintf("%s.*, %s.sess_id as session_id, %s.sess_time as session_timestamp, %s.id as session_id_number", tblUser, tblSess, tblSess, tblSess)
						join := fmt.Sprintf("JOIN %s on %s.id = %s.userid", tblSess, tblUser, tblSess)
						where := fmt.Sprintf("%s.id = ? AND %s.userid = ? AND %s.active = '1' AND %s.remove = '0' AND %s.block = '0' AND %s.suspend = '0'", tblSess, tblSess, tblUser, tblUser, tblUser, tblUser)

						err := config.DB.Table(tblUser).Select(sel).Joins(join).First(&user, where, verify.Data.SessionId, verify.Data.UserId).Error

						if err == nil {
							ctx.User = &user
						}
					}
					ctx.IsInternalServer = true
				}
			} else {
				ctx.IsApi = true
			}
			ctx.IsInternal = ctx.IsNative || ctx.IsWeb

			if xAppToken != "" {
				if _, err := config.FirebaseAppCheck.VerifyToken(xAppToken); err == nil {
					if client.Internal {
						ctx.IsInternal = true
					}
				}
			}
			if !ctx.IsNative && xDebug != "" {
				verify := util.VerifyToken[interface{}](xDebug, os.Getenv("AUTH_DEBUG_SECRET"), int64(time.Hour)*1)
				if verify.Verified {
					ctx.IsDebug = true
				}
			}
		}

		if options.Disable {
			ctx.Checklist = true
			c.Locals("ctx", ctx)
			return c.Next()
		}

		if clientId != "" {
			if err := config.DB.First(&client, "client_id = ?", clientId).Error; err != nil {
				fmt.Printf("Err Clients Database: %s\n\n", err.Error())
			}
		}

		if auth != "" && (ctx.IsNative || ctx.IsApi) {
			if clientId == "" {
				return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationMissingClientId)
			}
			var scope *[]string
			var grants, client_id string
			var token_database models.AccessToken
			auth_splice := strings.Split(auth, " ")
			auth_type := strings.ToLower(auth_splice[0])

			if auth_type == "bearer" {
				token_header := auth_splice[1]

				// Check Access Token
				if token_header == "" {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidAccessToken)
				}

				// Check Access Token in Database
				if err := config.DB.First(&token_database, "access_token = ?", token_header).Error; err != nil {
					fmt.Printf("Err Access Token Database: %s\n\n", err.Error())
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidAccessToken)
				}

				// Check AccessToken Client_ID && Access Token Expires
				if token_database.ClientId == nil || token_database.Expires == nil {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidAccessToken)
				}

				// Check Access Token Grants
				if ok := util.CheckGrants([]string{"client_credentials", "authorization_code"}, []string{*token_database.GrantTypes}); !ok {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidGrants)
				}

				// Check Access Token Client ID
				if *token_database.ClientId != clientId {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidClientId)
				}

				// Check Date
				date, err := time.Parse(time.RFC3339, *token_database.Expires)
				if err != nil {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationExpiredToken)
				}
				if time.Now().After(date) {
					return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationExpiredToken)
				}

				// Get Users
				if *token_database.GrantTypes == "authorization_code" && token_database.UserId != nil {
					var user models.UserContext
					var sel, join, where string

					var err error

					if ctx.IsNative {
						sel = fmt.Sprintf("%s.*, %s.sess_id as session_id, %s.sess_time as session_timestamp, %s.id as session_id_number", tblUser, tblSess, tblSess, tblSess)
						join = fmt.Sprintf("JOIN %s on %s.id = %s.userid", tblSess, tblUser, tblSess)
						where = fmt.Sprintf("%s.device_id = ? AND %s.userid = ? AND %s.active = '1' AND %s.remove = '0' AND %s.block = '0' AND %s.suspend = '0'", tblSess, tblSess, tblUser, tblUser, tblUser, tblUser)

						err = config.DB.Table(tblUser).Select(sel).Joins(join).First(&user, where, xDeviceId, *token_database.UserId).Error
					} else {
						sel = fmt.Sprintf("%s.*, %s.sess_id as session_id", tblUser, tblSess)

						join1 := fmt.Sprintf("JOIN %s on %s.id = %s.user_id", tblAccess, tblUser, tblAccess)
						join = fmt.Sprintf("JOIN %s on %s.user_id = %s.userid", tblSess, tblAccess, tblSess)

						where = fmt.Sprintf("%s.access_token = ? AND %s.user_id = ? AND %s.active = '1' AND %s.remove = '0' AND %s.block = '0' AND %s.suspend = '0'", tblAccess, tblAccess, tblUser, tblUser, tblUser, tblUser)

						err = config.DB.Table(tblUser).Select(sel).Joins(join1).Joins(join).First(&user, where, token_database.AccessToken, *token_database.UserId).Error
					}
					if err != nil {
						return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidAccessToken)
					}

					if !client.Publish && user.ID != client.UserId {
						if client.TestUserId != nil {
							if !util.ItemExists(*client.TestUserId, user.ID) {
								return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidClientIdDevelopment)
							}
						} else {
							return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationInvalidClientIdDevelopment)
						}
					}

					ctx.User = &user
				}
				scope = token_database.Scope
				grants = *token_database.GrantTypes
				client_id = client.ClientId
			} else {
				return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationNotSupported)
			}
			ctx.Client = &models.ClientContext{
				ClientId:    client_id,
				Scope:       scope,
				Grants:      grants,
				AccessToken: &token_database.AccessToken,
			}
		} else if pnAuth != "" {
			almostExpired := false

			verify := util.VerifyToken[AuthWeb](pnAuth, os.Getenv("AUTH_WEB_SECRET"), int64(time.Hour)*1)

			if !verify.Verified {
				http := fiber.StatusUnauthorized
				msg := "Missing token"
				if !verify.Info.Date {
					http = 440
					msg = "Token expired. Please refresh the browser"
				}
				return response.Authorization(http, response.ErrorAuthorizationCustom, &msg)
			}

			almostExpired = verify.Date.Add(time.Duration(time.Minute * 25)).Before(time.Now())
			ctx.AlmostExpired = almostExpired

			if xDebug != "" {
				verify := util.VerifyToken[AuthInternal](pnInternalPortalnesia, os.Getenv("AUTH_DEBUG_SECRET"), int64(time.Hour)*1)
				if verify.Verified {
					if verify.Data.UserId != nil && verify.Data.SessionId != nil {
						var user models.UserContext
						sel := fmt.Sprintf("%s.*, %s.sess_id as session_id, %s.sess_time as session_timestamp, %s.id as session_id_number", tblUser, tblSess, tblSess, tblSess)
						join := fmt.Sprintf("JOIN %s on %s.id = %s.userid", tblSess, tblUser, tblSess)
						where := fmt.Sprintf("%s.id = ? AND %s.userid = ? AND %s.active = '1' AND %s.remove = '0' AND %s.block = '0' AND %s.suspend = '0'", tblSess, tblSess, tblUser, tblUser, tblUser, tblUser)

						err := config.DB.Table(tblUser).Select(sel).Joins(join).First(&user, where, verify.Data.SessionId, verify.Data.UserId).Error

						if err == nil {
							ctx.User = &user
						}
					}
				}
				ctx.IsDebug = true
			} else if portalid != "" {
				verify := util.VerifyTokenAuth(portalid)
				if verify.Verified && verify.Data.Key != "" {
					key := verify.Data.Key
					auth, userid := key[0:60], key[60:]
					var user models.UserContext
					sel := fmt.Sprintf("%s.*, %s.sess_id as session_id, %s.sess_time as session_timestamp, %s.id as session_id_number", tblUser, tblSess, tblSess, tblSess)
					join := fmt.Sprintf("JOIN %s on %s.id = %s.userid", tblSess, tblUser, tblSess)
					where := fmt.Sprintf("%s.auth_key = ? AND %s.userid = ? AND %s.active = '1' AND %s.remove = '0' AND %s.block = '0' AND %s.suspend = '0'", tblSess, tblSess, tblUser, tblUser, tblUser, tblUser)

					err := config.DB.Table(tblUser).Select(sel).Joins(join).First(&user, where, auth, userid).Error

					if err == nil {
						ctx.User = &user
					}
				}
			}
		} else { // Native Without Authorization
			if !ctx.IsNative && !ctx.IsInternal && !ctx.IsInternalServer {
				return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationMissingAuthorization)
			}
			scope := []string{"superuser"}
			ctx.Client = &models.ClientContext{
				ClientId:    client.ClientId,
				Scope:       &scope,
				Grants:      "",
				AccessToken: nil,
			}
		}

		c.Locals("ctx", ctx)
		return c.Next()
	}
}

func OnlyInternal(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	if ctx.IsInternal || ctx.IsInternalServer {
		return c.Next()
	} else {
		return response.EndpointNotFound()
	}
}

func OnlyLogin(c *fiber.Ctx) error {
	ctx := c.Locals("ctx").(*models.Context)
	if ctx.User != nil {
		return c.Next()
	}
	if ctx.IsInternal {
		return response.Authorization(fiber.StatusForbidden, response.ErrorAuthorizationUnauthorizedLogin)
	} else {
		return response.Authorization(fiber.StatusForbidden, response.ErrorAuthorizationUnauthorizedGrants)
	}
}

func OnlySpecificScope(scope []string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := c.Locals("ctx").(*models.Context)
		if ctx.IsWeb || ctx.IsInternalServer {
			return c.Next()
		}
		if ctx.Client != nil {
			if ok := util.CheckScope(*ctx.Client.Scope, scope); ok {
				return c.Next()
			}
			return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationUnauthorizedScopes)
		}
		return response.Authorization(fiber.StatusUnauthorized, response.ErrorAuthorizationMissingAuthorization)
	}
}

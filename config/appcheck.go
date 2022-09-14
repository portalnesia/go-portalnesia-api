package config

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

const appCheckJWKSUrl = "https://firebaseappcheck.googleapis.com/v1beta/jwks"
const appCheckIssuer = "https://firebaseappcheck.googleapis.com/"

var (
	// ErrIncorrectAlgorithm is returned when the token is signed with a non-RSA256 algorithm.
	ErrAppCheckIncorrectAlgorithm = errors.New("token has incorrect algorithm")
	// ErrTokenType is returned when the token is not a JWT.
	ErrAppCheckTokenType = errors.New("token has incorrect type")
	// ErrTokenClaims is returned when the token claims cannot be decoded.
	ErrAppCheckTokenClaims = errors.New("token has incorrect claims")
	// ErrTokenAudience is returned when the token audience does not match the current project.
	ErrAppCheckTokenAudience = errors.New("token has incorrect audience")
	// ErrTokenIssuer is returned when the token issuer does not match Firebase's App Check service.
	ErrAppCheckTokenIssuer = errors.New("token has incorrect issuer")
	// ErrTokenSubject is returned when the token subject is empty or missing.
	ErrAppCheckTokenSubject = errors.New("token has empty or missing subject")
)

// VerifiedToken represents a verified App Check token.
type AppCheckVerifiedToken struct {
	Iss   string
	Sub   string
	Aud   []string
	Exp   time.Time
	Iat   time.Time
	AppID string
}

// Client is the interface for the Firebase App Check service.
type AppCheckClient struct {
	projectID string
	jwks      *keyfunc.JWKS
}

type AppCheckConfig struct {
	ProjectID string
	JWKSUrl   string
}

func NewAppCheck(ctx context.Context) (*AppCheckClient, error) {
	conf := &AppCheckConfig{
		ProjectID: "portalnesia",
		JWKSUrl:   appCheckJWKSUrl,
	}
	return appCheckNewClient(ctx, conf)
}

// NewClient creates a new App Check client.
func appCheckNewClient(ctx context.Context, conf *AppCheckConfig) (*AppCheckClient, error) {
	// TODO: Add support for overriding the HTTP client using the App one.
	jwks, err := keyfunc.Get(conf.JWKSUrl, keyfunc.Options{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &AppCheckClient{
		projectID: conf.ProjectID,
		jwks:      jwks,
	}, nil
}

func (c *AppCheckClient) VerifyToken(token string) (*AppCheckVerifiedToken, error) {
	// References for checks:
	// https://firebase.googleblog.com/2021/10/protecting-backends-with-app-check.html
	// https://github.com/firebase/firebase-admin-node/blob/master/src/app-check/token-verifier.ts#L106

	// The standard JWT parser also validates the expiration of the token
	// so we do not need dedicated code for that.
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if t.Header["alg"] != "RS256" {
			return nil, ErrAppCheckIncorrectAlgorithm
		}
		if t.Header["typ"] != "JWT" {
			return nil, ErrAppCheckTokenType
		}
		return c.jwks.Keyfunc(t)
	})
	if err != nil {
		return nil, err
	}

	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrAppCheckTokenClaims
	}

	rawAud := claims["aud"].([]interface{})
	aud := []string{}
	for _, v := range rawAud {
		aud = append(aud, v.(string))
	}

	if !contains(aud, "projects/"+c.projectID) {
		return nil, ErrAppCheckTokenAudience
	}

	// We check the prefix to make sure this token was issued
	// by the Firebase App Check service, but we do not check the
	// Project Number suffix because the Golang SDK only has project ID.
	//
	// This is consistent with the Firebase Admin Node SDK.
	if !strings.HasPrefix(claims["iss"].(string), appCheckIssuer) {
		return nil, ErrAppCheckTokenIssuer
	}

	if val, ok := claims["sub"].(string); !ok || val == "" {
		return nil, ErrAppCheckTokenSubject
	}

	return &AppCheckVerifiedToken{
		Iss:   claims["iss"].(string),
		Sub:   claims["sub"].(string),
		Aud:   aud,
		Exp:   time.Unix(int64(claims["exp"].(float64)), 0),
		Iat:   time.Unix(int64(claims["iat"].(float64)), 0),
		AppID: claims["sub"].(string),
	}, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

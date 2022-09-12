package webauthn

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/portalnesia/go-utils/goment"
	"portalnesia.com/api/models"
	util "portalnesia.com/api/utils"
)

var config = webauthn.Config{
	RPDisplayName: "Portalnesia",
	RPID:          "portalnesia.com",
	RPIcon:        util.StaticUrl("icon/PN-Logo-512.png"),
	RPOrigin:      "https://portalnesia.com",
}

type WebauthnUser struct {
	models.User
	Credentials []webauthn.Credential
}

func (u *WebauthnUser) WebAuthnID() []byte {
	user_id := strconv.Itoa(int(u.ID))
	return []byte(user_id)
}

func (u *WebauthnUser) WebAuthnName() string {
	return u.Username
}

func (u *WebauthnUser) WebAuthnDisplayName() string {
	return u.Username
}

func (u *WebauthnUser) WebAuthnIcon() string {
	return util.StaticUrl("icon/PN-Logo-512.png")
}

func (u *WebauthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func newUser(u models.User) webauthn.User {
	credentials := make([]webauthn.Credential, len(u.Webauthn.WebauthnKeys))

	for i, c := range u.Webauthn.WebauthnKeys {
		credentials[i] = webauthn.Credential{
			ID:        c.ID,
			PublicKey: []byte(c.PublicKey),
		}
	}

	user := WebauthnUser{
		User:        u,
		Credentials: credentials,
	}

	return &user
}

type TokenRequest[T any] struct {
	Token   string `json:"token"`
	Date    string `json:"date"`
	Request *T     `json:"challenge"`
}

type AuthResponse[T any, Y any] struct {
	Token   string `json:"token"`
	Request T      `json:"payload"`
	Parsed  Y      `json:"-"`
}

func BeginLogin(u models.User) (*protocol.CredentialAssertion, string, error) {
	authn, err := webauthn.New(&config)
	if err != nil {
		return nil, "", err
	}

	user := newUser(u)

	creds, _, err := authn.BeginLogin(user)
	if err != nil {
		return nil, "", err
	}

	date, _ := goment.New()
	token := util.CreateToken(TokenRequest[protocol.CredentialAssertion]{
		Request: creds,
		Token:   os.Getenv("WEBAUTHN_LOGIN_SECRET"),
		Date:    date.PNformat(),
	})
	return creds, token, nil
}

func ParseLoginRequestBody(c *fiber.Ctx) (*AuthResponse[protocol.CredentialAssertionResponse, protocol.ParsedCredentialAssertionData], error) {
	var resp AuthResponse[protocol.CredentialAssertionResponse, protocol.ParsedCredentialAssertionData]
	err := c.BodyParser(&resp)
	if err != nil {
		return nil, protocol.ErrBadRequest.WithDetails("Invalid Request")
	}

	car := resp.Request
	if car.ID == "" {
		return nil, protocol.ErrBadRequest.WithDetails("CredentialAssertionResponse with ID missing")
	}
	_, err = base64.RawURLEncoding.DecodeString(car.ID)
	if err != nil {
		return nil, protocol.ErrBadRequest.WithDetails("CredentialAssertionResponse with ID not base64url encoded")
	}
	if car.Type != "public-key" {
		return nil, protocol.ErrBadRequest.WithDetails("CredentialAssertionResponse with bad type")
	}

	var par protocol.ParsedCredentialAssertionData
	par.ID, par.RawID, par.Type, par.ClientExtensionResults = car.ID, car.RawID, car.Type, car.ClientExtensionResults
	par.Raw = car

	par.Response.Signature = car.AssertionResponse.Signature
	par.Response.UserHandle = car.AssertionResponse.UserHandle
	err = json.Unmarshal(car.AssertionResponse.ClientDataJSON, &par.Response.CollectedClientData)
	if err != nil {
		return nil, err
	}

	err = par.Response.AuthenticatorData.Unmarshal(car.AssertionResponse.AuthenticatorData)
	if err != nil {
		return nil, protocol.ErrParsingData.WithDetails("Error unmarshalling auth data")
	}
	resp.Parsed = par
	return &resp, nil
}

func Login(resp *AuthResponse[protocol.CredentialAssertionResponse, protocol.ParsedCredentialAssertionData], u models.User) (*webauthn.Credential, error) {
	token := util.VerifyToken[TokenRequest[protocol.CredentialAssertion]](resp.Token, os.Getenv("WEBAUTHN_LOGIN_SECRET"), int64(time.Minute)*15)

	authn, err := webauthn.New(&config)
	if err != nil {
		return nil, err
	}

	user := newUser(u)
	creds := user.WebAuthnCredentials()
	var allowedCredentials = make([][]byte, len(creds))

	for i, credential := range creds {
		allowedCredentials[i] = credential.ID
	}

	sess := webauthn.SessionData{
		Challenge:            token.Data.Request.Response.Challenge.String(),
		UserID:               user.WebAuthnID(),
		AllowedCredentialIDs: allowedCredentials,
		UserVerification:     authn.Config.AuthenticatorSelection.UserVerification,
		Extensions:           protocol.AuthenticationExtensions{},
	}

	return authn.ValidateLogin(user, sess, &resp.Parsed)
}

func BeginRegister(u models.User) (*protocol.CredentialCreation, string, error) {
	authn, err := webauthn.New(&config)
	if err != nil {
		return nil, "", err
	}

	user := newUser(u)

	creds, _, err := authn.BeginRegistration(user)
	if err != nil {
		return nil, "", err
	}

	date, _ := goment.New()
	token := util.CreateToken(TokenRequest[protocol.CredentialCreation]{
		Request: creds,
		Token:   os.Getenv("WEBAUTHN_REGISTER_SECRET"),
		Date:    date.PNformat(),
	})
	return creds, token, nil
}

func ParseRegisterRequestBody(c *fiber.Ctx) (*AuthResponse[protocol.CredentialCreationResponse, protocol.ParsedCredentialCreationData], error) {
	var resp AuthResponse[protocol.CredentialCreationResponse, protocol.ParsedCredentialCreationData]
	err := c.BodyParser(&resp)
	if err != nil {
		return nil, protocol.ErrBadRequest.WithDetails("Invalid Request")
	}

	ccr := resp.Request
	if ccr.ID == "" {
		return nil, protocol.ErrBadRequest.WithDetails("Parse error for Registration").WithInfo("Missing ID")
	}

	testB64, err := base64.RawURLEncoding.DecodeString(ccr.ID)
	if err != nil || !(len(testB64) > 0) {
		return nil, protocol.ErrBadRequest.WithDetails("Parse error for Registration").WithInfo("ID not base64.RawURLEncoded")
	}

	if ccr.PublicKeyCredential.Credential.Type == "" {
		return nil, protocol.ErrBadRequest.WithDetails("Parse error for Registration").WithInfo("Missing type")
	}

	if ccr.PublicKeyCredential.Credential.Type != "public-key" {
		return nil, protocol.ErrBadRequest.WithDetails("Parse error for Registration").WithInfo("Type not public-key")
	}

	var pcc protocol.ParsedCredentialCreationData
	pcc.ID, pcc.RawID, pcc.Type, pcc.ClientExtensionResults = ccr.ID, ccr.RawID, ccr.Type, ccr.ClientExtensionResults
	pcc.Raw = ccr

	parsedAttestationResponse, err := ccr.AttestationResponse.Parse()
	if err != nil {
		return nil, protocol.ErrParsingData.WithDetails("Error parsing attestation response")
	}

	pcc.Response = *parsedAttestationResponse

	resp.Parsed = pcc
	return &resp, nil
}

func Register(resp *AuthResponse[protocol.CredentialCreationResponse, protocol.ParsedCredentialCreationData], u models.User) (*webauthn.Credential, error) {
	token := util.VerifyToken[TokenRequest[protocol.CredentialAssertion]](resp.Token, os.Getenv("WEBAUTHN_REGISTER_SECRET"), int64(time.Minute)*15)

	authn, err := webauthn.New(&config)
	if err != nil {
		return nil, err
	}

	user := newUser(u)
	creds := user.WebAuthnCredentials()
	var allowedCredentials = make([][]byte, len(creds))

	for i, credential := range creds {
		allowedCredentials[i] = credential.ID
	}

	sess := webauthn.SessionData{
		Challenge:            token.Data.Request.Response.Challenge.String(),
		UserID:               user.WebAuthnID(),
		AllowedCredentialIDs: allowedCredentials,
		UserVerification:     authn.Config.AuthenticatorSelection.UserVerification,
		Extensions:           protocol.AuthenticationExtensions{},
	}

	return authn.CreateCredential(user, sess, &resp.Parsed)
}

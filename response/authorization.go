package response

const (
	ErrorAuthorizationCustom                     = 100
	ErrorAuthorizationMissingClientId            = 111
	ErrorAuthorizationMissingAccessToken         = 112
	ErrorAuthorizationMissingAuthorization       = 113
	ErrorAuthorizationExpiredToken               = 121
	ErrorAuthorizationExpiredParams              = 122
	ErrorAuthorizationExpiredSignature           = 123
	ErrorAuthorizationInvalidClientId            = 131
	ErrorAuthorizationInvalidAccessToken         = 132
	ErrorAuthorizationInvalidCredentials         = 133
	ErrorAuthorizationInvalidGrants              = 134
	ErrorAuthorizationInvalidSignature           = 135
	ErrorAuthorizationInvalidClientIdDevelopment = 136
	ErrorAuthorizationUnauthorizedGrants         = 141
	ErrorAuthorizationUnauthorizedScopes         = 142
	ErrorAuthorizationUnauthorizedUser           = 143
	ErrorAuthorizationUnauthorizedLogin          = 144
	ErrorAuthorizationNotSupported               = 140
	ErrorAuthorizationJwt                        = 150
	ErrorAuthorizationBlockClientId              = 161
)

func getAuthorizationMessage(err int) string {
	switch err {
	case ErrorAuthorizationMissingClientId:
		return "No `PN-Client-Id` provided"
	case ErrorAuthorizationMissingAccessToken:
		return "No `access token` provided"
	case ErrorAuthorizationMissingAuthorization:
		return "Missing authorization token"
	case ErrorAuthorizationExpiredToken:
		return "The `access token` provided has expired"
	case ErrorAuthorizationExpiredParams:
		return ""
	case ErrorAuthorizationExpiredSignature:
		return "The `signature` provided has expired"
	case ErrorAuthorizationInvalidClientId:
		return "The `client id` provided is invalid"
	case ErrorAuthorizationInvalidAccessToken:
		return "The `access token` provided is invalid"
	case ErrorAuthorizationInvalidCredentials:
		return "The `credentials` provided is invalid"
	case ErrorAuthorizationInvalidGrants:
		return "The `grant type` is unauthorized for this client id"
	case ErrorAuthorizationInvalidSignature:
		return "The `signature` provided is invalid"
	case ErrorAuthorizationInvalidClientIdDevelopment:
		return "The `client id` provided is invalid or still in development mode"
	case ErrorAuthorizationUnauthorizedGrants:
		return "The `grant type` is unauthorized for this request"
	case ErrorAuthorizationUnauthorizedScopes:
		return "The `scope` is unauthorized for this request"
	case ErrorAuthorizationUnauthorizedUser:
		return "The `username` is unauthorized with this `access_token`"
	case ErrorAuthorizationUnauthorizedLogin:
		return "You must login to make this request."
	case ErrorAuthorizationNotSupported:
		return "`Authorization` not supported"
	case ErrorAuthorizationJwt:
		return "The `JWT` provided is invalid"
	case ErrorAuthorizationBlockClientId:
		return "The `client id` provided is blocked"
	}
	return "Unknown"
}

func Authorization(http int, err int, msg ...*string) *Error {
	var ms string
	if msg != nil && msg[0] != nil {
		ms = *msg[0]
	} else {
		ms = getAuthorizationMessage(err)
	}
	return NewError(http, err, "authorization", ms)
}

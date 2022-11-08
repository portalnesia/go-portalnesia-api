package models

type Context struct {
	// IS DEVELOPER CLIENT (EXTERNAL APPLICATION WITH OAUTH2)
	IsApi bool
	// IS WEB APPLICATION
	IsWeb bool
	// IS INTERNAL PORTALNESIA
	IsInternal bool
	// IS ACCESSED FROM LOCALHOST PORTALNESIA
	IsDebug bool
	// IS USE NATIVE APPLICATION
	IsNative bool
	// IS ACCESSED FROM PHP OR INTERNAL SERVER
	IsInternalServer bool
	// CHECKLIST FOR FIRST AUTHORIZATION AND SECOND
	Checklist     bool
	AlmostExpired bool
	
}

var CtxDefaultValue = Context{
	IsApi:            false,
	IsWeb:            false,
	IsInternal:       false,
	IsDebug:          false,
	IsNative:         false,
	IsInternalServer: false,
	Checklist:        false,
	AlmostExpired:    false,
}

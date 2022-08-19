package util

import (
	"fmt"
	"net/url"
	"os"

	"github.com/portalnesia/go-utils"
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

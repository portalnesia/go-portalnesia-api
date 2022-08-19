package util

import (
	"fmt"
	"os"
	"testing"
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

package utils

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/portalnesia/go-utils"
)

const defaultWidth int = 400

func isUnsplashImage(url string) bool {
	regex := regexp.MustCompile(`https?\:\/\/(www.)?unsplash\.com`)
	return regex.MatchString(url)
}

func parseImage(s *goquery.Selection, ads bool) *goquery.Selection {
	var regex *regexp.Regexp
	src_string := ""
	src := ""
	withPng := false
	src_string = s.AttrOr("src", "")
	if src_string == "" {
		src_string = s.AttrOr("data-src", "")
	}

	attr, exist := s.Attr("data-png")
	if exist {
		if attr == "true" {
			withPng = true
		}
	}

	domain, _ := url.Parse(src_string)
	src_domain := domain.Host
	isUnsplash := isUnsplashImage(src_string)

	// SRC
	if src_domain == "content.portalnesia.com" {
		src = src_string
	} else if isUnsplash {
		regex = regexp.MustCompile(`(\?|\&)w\=\d+`)
		r := regex.ReplaceAllString(src_string, "")
		src = fmt.Sprintf("%s&auto=compress&mark=%s&mark-scale=5&mark-align=middle", r, url.QueryEscape(StaticUrl("watermark.png")))
	} else {
		src = AnalyzeStaticUrl(src_string)
	}

	// SIZE
	attr, exist = s.Attr("width")
	width := defaultWidth
	if exist {
		c, err := strconv.Atoi(attr)
		if err == nil {
			width = c
		}
	}

	// CAPTION
	caption := ""
	elem := s.ParentsFiltered("figure").Find("figcaption")
	if elem.Length() > 0 {
		h, _ := elem.Html()
		caption = h
	} else {
		if s.AttrOr("alt", "") != "" {
			caption = s.AttrOr("alt", "")
		}
	}
	caption = utils.Clean(caption)
	caption = regexp.MustCompile(`\s+https?\:\/\/\S+`).ReplaceAllString(caption, "")
	caption = regexp.MustCompile(`\n`).ReplaceAllString(caption, "")

	regex = regexp.MustCompile(`(\?|\&)image\=pixabay`)
	if src_domain == "content.portalnesia.com" && regex.MatchString(src) {
		link, _ := url.Parse(src)
		if link.Query().Has("image") {
			imageUrl, _ := url.QueryUnescape(link.Query().Get("image"))
			if path.Ext(imageUrl) == ".png" {
				withPng = true
			}
		}
		if link.Query().Has("output") && link.Query().Get("output") == "png" {
			withPng = true
		}
	}

	if withPng {
		s.SetAttr("data-png", "true")
	}
	s.SetAttr("loading", "lazy")
	s.AddClass("image-container")
	s.AddClass("loading")
	s.RemoveAttr("width")
	s.RemoveAttr("height")

	if caption != "" {
		s.SetAttr("alt", caption)
	}

	if !ads {
		if isUnsplash {
			s.SetAttr("src", fmt.Sprintf("%s&w=800&q=80&mark-pad=10", src))
		} else {
			s.SetAttr("src", fmt.Sprintf("%s&size=%d", src, width))
		}
		s.RemoveAttr("data-src")
		return s.Clone()
	} else {
		s.RemoveAttr("src")
		if isUnsplash {
			s.SetAttr("data-src", fmt.Sprintf("%s&w=300&q=55&mark-pad=5", src))
		} else {
			s.SetAttr("data-src", fmt.Sprintf("%s&size=%d", src, width))
		}

		r := strings.NewReader("<a></a>")
		aDoc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			log.Fatalln(err)
		}

		a := aDoc.FindMatcher(goquery.Single("a"))
		a.SetAttr("data-fancybox", "true")
		if isUnsplash {
			a.SetAttr("data-src", fmt.Sprintf("%s&w=800&q=80&mark-pad=10", src))
		} else {
			a.SetAttr("data-src", src)
		}

		if caption != "" {
			a.SetAttr("data-caption", caption)
		}
		newImg := s.Clone()
		a.AppendSelection(newImg)
		return a
	}
}

func BlogEncode(html string, ads bool) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		log.Fatalln(err)
	}

	var newElem *goquery.Selection
	picture := doc.Find("picture")
	imgTag := doc.Find("img")
	figure := doc.Find("figure")
	pTag := doc.Find("p")

	picture.Each(func(_ int, p *goquery.Selection) {
		p.Find("img").Each(func(i int, s *goquery.Selection) {
			newElem = parseImage(s, ads)
			p.ReplaceWithSelection(newElem)
		})
	})

	imgTag.Each(func(i int, s *goquery.Selection) {
		newElem = parseImage(s, ads)
		s.ReplaceWithSelection(newElem)
	})

	figure.Each(func(i int, f *goquery.Selection) {
		if f.HasClass("table") {
			f.Find(".table").Each(func(i int, s *goquery.Selection) {
				f.ReplaceWithSelection(s)
			})
		} else {
			style := f.AttrOr("style", "")
			if regexp.MustCompile(`width`).MatchString(style) {
				f.SetAttr("style", "max-width:400px;width:90%")
			}
		}
	})

	pTag.Each(func(i int, p *goquery.Selection) {
		if i == int(math.Round(float64(pTag.Length())/3)) {
			p.PrependHtml(`<div data-portalnesia-action="ads" data-ads="300"></div>`)
		}
		if i == int(math.Round(2*float64(pTag.Length())/3)) {
			p.PrependHtml(`<div data-portalnesia-action="ads" data-ads="468"></div>`)
		}
	})

	out, err := doc.Html()
	if err != nil {
		log.Fatalln(err)
	}
	return out
}

func NewsEncode(html string) string {
	r := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		log.Fatalln(err)
	}

	str, _ := doc.Html()

	str = strings.ReplaceAll(str, "[IKLAN_IKLAN]</p>", "</p><div data-portalnesia-action=\"ads\" data-ads=\"300\"></div>")
	str = strings.ReplaceAll(str, "[IKLAN_IKLAN_IKLAN]</p>", "</p><div data-portalnesia-action=\"ads\" data-ads=\"468\"></div>")

	return str
}

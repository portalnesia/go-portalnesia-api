package utils

import (
	"fmt"
	"os"
)

type Path struct {
	content string
	php     string
}

func NewPath() *Path {
	content := os.Getenv("PORTALNESIA_CONTENT_ROOT")
	php := os.Getenv("PORTALNESIA_PHP_ROOT")
	return &Path{
		content: content,
		php:     php,
	}
}

func getArgs(p ...string) string {
	var path string
	if len(p) == 1 {
		path = fmt.Sprintf("/%s", p[0])
	}
	return path
}

func (p *Path) Content(path ...string) string {
	pth := getArgs(path...)
	return fmt.Sprintf("%s%s", p.content, pth)
}

func (p *Path) PHP(path ...string) string {
	pth := getArgs(path...)
	return fmt.Sprintf("%s%s", p.php, pth)
}

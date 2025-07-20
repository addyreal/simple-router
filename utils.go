package router

import (
	"path"
	"strings"
)

func split(a string) []string {
	b := path.Clean("/" + a)
	c := strings.Trim(b, "/")
	return strings.Split(c, "/")
}

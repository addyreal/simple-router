package url

import (
	"path"
	"strings"
)

func Split(a string) []string {
	b := path.Clean("/" + a)
	c := strings.Trim(b, "/")
	return strings.Split(c, "/")
}

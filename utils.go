package router

import (
	"path"
	"net/http"
	"strings"
)

func split(a string) []string {
	b := path.Clean("/" + a)
	c := strings.Trim(b, "/")
	return strings.Split(c, "/")
}

func compose(a, b func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	if a == nil {
		a = func(h http.HandlerFunc) http.HandlerFunc {return h}
	}

	return func(aC, bC func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
		return func(h http.HandlerFunc) http.HandlerFunc {
			return bC(aC(h))
		}
	}(a, b)
}

func wrap(a func(http.HandlerFunc) http.HandlerFunc, b http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a(b)(w, r)
	}
}

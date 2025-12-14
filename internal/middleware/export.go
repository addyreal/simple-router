package middleware

import "net/http"

func Identity(h http.HandlerFunc) http.HandlerFunc {
	return h
}

func Compose(a, b func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	if a == nil {
		a = Identity
	}
	if b == nil {
		b = Identity
	}

	return func(aC, bC func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
		return func(h http.HandlerFunc) http.HandlerFunc {
			return aC(bC(h))
		}
	}(a, b)
}

func Wrap(a func(http.HandlerFunc) http.HandlerFunc, b http.HandlerFunc) http.HandlerFunc {
	if a == nil {
		a = Identity
	}
	if b == nil {
		panic("Tried to wrap a nil handler")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		a(b)(w, r)
	}
}

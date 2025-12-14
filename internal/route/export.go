package route

import "net/http"

type Route struct {
	Class   int
	Handler http.HandlerFunc
}

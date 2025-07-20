package router

import (
	"net/http"
)

func Init() *temp {
	return &temp {
		headers: nil,
		middleware: func(h http.HandlerFunc) http.HandlerFunc {return h},
		notfound: nil,
		recovery: nil,
		gets: make(map[string]http.HandlerFunc),
		posts: make(map[string]http.HandlerFunc),
	}
}

func (x *temp) SetHeaders(y func(http.ResponseWriter)) {
	x.headers = y
}

func (x *temp) AppendMiddleware(y func(http.HandlerFunc) http.HandlerFunc) {
	a := x.middleware
	x.middleware = func(b http.HandlerFunc) http.HandlerFunc {
		return a(b)
	}
}

func (x *temp) SetNotFound(y http.HandlerFunc) {
	x.notfound = y
}

func (x *temp) SetRecovery(y http.HandlerFunc) {
	x.recovery = y
}

func (x *temp) AddGet(y string, z http.HandlerFunc) {
	x.gets[y] = z
}

func (x *temp) AddPost(y string, z http.HandlerFunc) {
	x.posts[y] = z
}

func (x *temp) Get() http.HandlerFunc {
	if x.headers == nil || x.notfound == nil || x.recovery == nil {
		panic("Router unimplemented")
	}

	g := buildTrie(x.gets)
	p := buildTrie(x.posts)

	return x.middleware(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				x.recovery(w, r)
			}
		}()

		x.headers(w)

		var which *node
		switch r.Method {
			case http.MethodGet:
				which = g
			case http.MethodPost:
				which = p
			default:
				x.notfound(w, r)
				return
		}

		parts := split(r.URL.Path)
		handler := which.walk(parts)
		if handler == nil {
			x.notfound(w, r)
			return
		}

		handler(w, r)
	})
}

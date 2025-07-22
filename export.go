package router

import (
	"net/http"
)

func Init() *temp {
	return &temp {
		headers: nil,
		notfound: nil,
		recovery: nil,
		gets: nil,
		posts: nil,
		middleware: make(map[int]func(http.HandlerFunc) http.HandlerFunc),
	}
}

func (x *temp) SetHeaders(y func(http.ResponseWriter)) {
	x.headers = y
}

func (x *temp) SetNotFound(y http.HandlerFunc) {
	x.notfound = y
}

func (x *temp) SetRecovery(y http.HandlerFunc) {
	x.recovery = y
}

func (x *temp) AddMiddleware(y func(http.HandlerFunc) http.HandlerFunc) {
	x.middleware[-1] = compose(x.middleware[-1], y)
}

func (x *temp) AddGet(g uint8, y string, z http.HandlerFunc) {
	x.gets = append(x.gets, route{group: g, path: y, handler: z})
}

func (x *temp) AddPost(g uint8, y string, z http.HandlerFunc) {
	x.posts = append(x.posts, route{group: g, path: y, handler: z})
}

func (x *temp) AppendMiddleware(gs []uint, y func(http.HandlerFunc) http.HandlerFunc) {
	for _, g := range gs {
		x.middleware[int(g)] = compose(x.middleware[int(g)], y)
	}
}

func (x *temp) Get() http.HandlerFunc {
	if x.headers == nil || x.notfound == nil || x.recovery == nil {
		panic("Router unimplemented")
	}

	if x.middleware[-1] == nil {
		x.middleware[-1] = func(h http.HandlerFunc) http.HandlerFunc {return h}
	}

	a := buildTrie(x.gets, x.middleware)
	b := buildTrie(x.posts, x.middleware)

	return func(w http.ResponseWriter, r *http.Request) {
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
				which = a
			case http.MethodPost:
				which = b
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
	}
}

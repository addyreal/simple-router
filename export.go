package router

import (
	"net/http"
)

func HeadFromGet(x http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &resrec{ResponseWriter: w}
		x(rec, r)
	}
}

func Init() *temp {
	return &temp {
		headers: nil,
		notfound: nil,
		recovery: nil,
		heads: nil,
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

func (x *temp) SetRecovery(y func(any, http.ResponseWriter, *http.Request)) {
	x.recovery = y
}

func (x *temp) AddMiddleware(y func(http.HandlerFunc) http.HandlerFunc) {
	x.middleware[-1] = compose(x.middleware[-1], y)
}

func (x *temp) AddHead(g uint8, y string, z http.HandlerFunc) {
	x.heads = append(x.heads, route{group: g, path: y, handler: z})
}

func (x *temp) AddGet(g uint8, y string, z http.HandlerFunc) {
	x.gets = append(x.gets, route{group: g, path: y, handler: z})
}

func (x *temp) AddPost(g uint8, y string, z http.HandlerFunc) {
	x.posts = append(x.posts, route{group: g, path: y, handler: z})
}

func (x *temp) AppendMiddleware(gs []uint8, y func(http.HandlerFunc) http.HandlerFunc) {
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

	a := buildTrie(x.heads, x.middleware)
	b := buildTrie(x.gets, x.middleware)
	c := buildTrie(x.posts, x.middleware)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				x.recovery(err, w, r)
			}
		}()

		x.headers(w)

		var which *node
		switch r.Method {
			case http.MethodHead:
				which = a
			case http.MethodGet:
				which = b
			case http.MethodPost:
				which = c
			default:
				x.notfound(w, r)
				return
		}

		parts := split(r.URL.Path)
		handler := which.walk(parts)
		if handler == nil {
			if r.Method == http.MethodHead {
				HeadFromGet(x.notfound).ServeHTTP(w, r)
				return
			} else {
				x.notfound(w, r)
				return
			}
		}

		handler(w, r)
	}
}

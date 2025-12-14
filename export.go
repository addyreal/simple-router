package router

import (
	"github.com/addyreal/simple-router/internal/head"
	"github.com/addyreal/simple-router/internal/middleware"
	"github.com/addyreal/simple-router/internal/route"
	"github.com/addyreal/simple-router/internal/trie"
	"net/http"
)

func HeadOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w2 := &head.HeaderWriter{ResponseWriter: w}
		h(w2, r)
	}
}

type temp struct {
	notfound   http.HandlerFunc
	recovery   func(any, http.ResponseWriter, *http.Request)
	routes     map[string]map[string]route.Route
	middleware map[int]func(http.HandlerFunc) http.HandlerFunc
}

func Init() *temp {
	return &temp{
		notfound:   nil,
		recovery:   nil,
		routes:     make(map[string]map[string]route.Route),
		middleware: make(map[int]func(http.HandlerFunc) http.HandlerFunc),
	}
}

func (x *temp) SetNotFound(h http.HandlerFunc) {
	x.notfound = h
}

func (x *temp) SetRecovery(h func(any, http.ResponseWriter, *http.Request)) {
	x.recovery = h
}

func (x *temp) AddMiddleware(m func(http.HandlerFunc) http.HandlerFunc) {
	x.middleware[-1] = middleware.Compose(x.middleware[-1], m)
}

func (x *temp) AppendMiddleware(c int, m func(http.HandlerFunc) http.HandlerFunc) {
	if c < 0 {
		panic("Negative groups are reserved")
	}
	x.middleware[c] = middleware.Compose(x.middleware[c], m)
}

func (x *temp) Add(m string, c int, p string, h http.HandlerFunc) {
	if c < 0 {
		panic("Negative groups are reserved")
	}
	if x.routes[m] == nil {
		x.routes[m] = make(map[string]route.Route)
	}
	x.routes[m][p] = route.Route{Class: c, Handler: h}
}

func (x *temp) Get() http.HandlerFunc {
	if x.notfound == nil || x.recovery == nil {
		panic("Router unimplemented")
	}
	if x.middleware[-1] == nil {
		x.middleware[-1] = middleware.Identity
	}

	tries := trie.BuildTries(x.routes, x.middleware)
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				x.recovery(err, w, r)
				return
			}
		}()

		arg := r.Method
		if arg == "" {
			arg = "GET"
		}

		tree, ok := tries[arg]
		if ok == false {
			if r.Method == http.MethodHead {
				HeadOnly(x.notfound)(w, r)
				return
			} else {
				x.notfound(w, r)
				return
			}
		}

		handler := tree.Walk(r.URL.Path)
		if handler == nil {
			if r.Method == http.MethodHead {
				HeadOnly(x.notfound)(w, r)
				return
			} else {
				x.notfound(w, r)
				return
			}
		}

		handler(w, r)
	}
}

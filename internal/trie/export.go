package trie

import (
	"github.com/addyreal/simple-router/internal/middleware"
	"github.com/addyreal/simple-router/internal/route"
	"github.com/addyreal/simple-router/internal/url"
	"net/http"
)

type node struct {
	children map[string]*node
	handler  http.HandlerFunc
}

func makeNode() *node {
	return &node{
		children: make(map[string]*node),
		handler:  nil,
	}
}

func (x *node) addPath(p string, h http.HandlerFunc) {
	for _, part := range url.Split(p) {
		if x.children[part] == nil {
			x.children[part] = makeNode()
		}

		x = x.children[part]
	}

	x.handler = h
}

func build(a map[string]route.Route, b map[int]func(http.HandlerFunc) http.HandlerFunc) *node {
	res := makeNode()
	for k, v := range a {
		if b[v.Class] == nil {
			b[v.Class] = middleware.Identity
		}

		res.addPath(k, middleware.Wrap(middleware.Compose(b[-1], b[v.Class]), v.Handler))
	}

	return res
}

func BuildTries(a map[string]map[string]route.Route, b map[int]func(http.HandlerFunc) http.HandlerFunc) map[string]*node {
	res := make(map[string]*node)
	for k, v := range a {
		res[k] = build(v, b)
	}

	return res
}

func (x *node) Walk(p string) http.HandlerFunc {
	for _, part := range url.Split(p) {
		if x.children[part] == nil {
			return nil
		}

		x = x.children[part]
	}

	return x.handler
}

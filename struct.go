package router

import (
	"net/http"
)

type route struct {
	group		uint8
	path		string
	handler		http.HandlerFunc
}

type temp struct {
	headers		func(http.ResponseWriter)
	notfound	http.HandlerFunc
	recovery	http.HandlerFunc
	gets		[]route
	posts		[]route
	middleware	map[int]func(http.HandlerFunc) http.HandlerFunc
}

type node struct {
	children	map[string]*node
	handler		http.HandlerFunc
}

func makeNode() *node {
	return &node{children: make(map[string]*node)}
}

func (x *node) add(a string, b http.HandlerFunc) {
	parts := split(a)

	for _, part := range parts {
		if x.children[part] == nil {
			x.children[part] = makeNode()
		}

		x = x.children[part]
	}
	x.handler = b
}

func (x *node) walk(parts []string) http.HandlerFunc {
	for _, part := range parts {
		a := x.children[part]
		if a == nil {
			return nil
		}

		x = a
	}

	return x.handler
}

func buildTrie(a []route, mw map[int]func(http.HandlerFunc) http.HandlerFunc) *node {
	trie := makeNode()

	for _, r := range a {
		g := mw[-1]
		m := mw[int(r.group)]
		h := r.handler
		if m != nil {
			m = compose(g, m)
		}
		h = wrap(m, h)
		trie.add(r.path, h)
	}

	return trie
}

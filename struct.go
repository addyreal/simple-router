package router

import (
	"net/http"
)

type temp struct {
	headers		func(http.ResponseWriter)
	middleware	func(http.HandlerFunc) http.HandlerFunc
	notfound	http.HandlerFunc
	recovery	http.HandlerFunc
	gets		map[string]http.HandlerFunc
	posts		map[string]http.HandlerFunc
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

func buildTrie(a map[string]http.HandlerFunc) *node {
	trie := makeNode()
	for b, c := range a {
		trie.add(b, c)
	}

	return trie
}

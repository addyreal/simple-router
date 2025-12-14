#!/bin/bash

DIR="$(dirname "${BASH_SOURCE[0]}")"
TESTFILE="router_test.go"
[[ -w "$DIR" ]] || { echo "Error creating file $TESTFILE"; exit 1; }
[[ -e $TESTFILE ]] && { echo "Error creating file $TESTFILE"; exit 1; }
which go >/dev/null 2>&1 || { echo "Error invoking go"; exit 1; }
cd "$DIR"
tee $TESTFILE >/dev/null <<'EOF'
package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func _Ok(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OKK!!!"))
}
func _NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404!!!"))
}
func _Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panicking.")
}
func _Recover(err any, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	w.Write([]byte(err.(string)))
}
func _MiddlewareContextSetFirst(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "key", "first")
		n.ServeHTTP(w, r.WithContext(ctx))
	}
}
func _MiddlewareContextSetSecond(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value("key").(string)
		ctx := context.WithValue(r.Context(), "key", val+"second")
		n.ServeHTTP(w, r.WithContext(ctx))
	}
}
func _HandlerWriteContext(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Context().Value("key").(string)))
}
func _MiddlewareGlobalContextSetOne(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "globalkey", "one")
		n.ServeHTTP(w, r.WithContext(ctx))
	}
}
func _MiddlewareGlobalContextSetTwo(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "globalkey", "two")
		n.ServeHTTP(w, r.WithContext(ctx))
	}
}
func _HandlerWriteGlobalContext(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.Context().Value("globalkey").(string)))
}

func TestRouter(t *testing.T) {
	b := Init()
	b.SetNotFound(_NotFound)
	b.SetRecovery(_Recover)
	b.AddMiddleware(_MiddlewareGlobalContextSetOne)

	b.Add("GET", 0, "/", _Ok)

	b.Add("GET", 0, "/panic", _Panic)

	b.Add("GET", 0, "/globalctx", _HandlerWriteGlobalContext)

	b.AppendMiddleware(1, _MiddlewareGlobalContextSetTwo)
	b.Add("GET", 1, "/globalafter", _HandlerWriteGlobalContext)

	b.AppendMiddleware(2, _MiddlewareContextSetFirst)
	b.Add("GET", 2, "/singlectx", _HandlerWriteContext)

	b.AppendMiddleware(3, _MiddlewareContextSetFirst)
	b.AppendMiddleware(3, _MiddlewareContextSetSecond)
	b.Add("GET", 3, "/doublectx", _HandlerWriteContext)

	router := b.Get()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{"Ok", "GET", "/", "OKK!!!", http.StatusOK},
		{"Ok no delim", "GET", "", "OKK!!!", http.StatusOK},
		{"Empty method", "", "/", "OKK!!!", http.StatusOK},
		{"Not found", "GET", "/AAAA", "404!!!", http.StatusNotFound},
		{"Head not found", "HEAD", "/AAAA", "", http.StatusNotFound},
		{"Wrong method", "POST", "/", "404!!!", http.StatusNotFound},
		{"Panic recovery", "GET", "/panic", "Panicking.", 201},
		{"Global middleware", "GET", "/globalctx", "one", http.StatusOK},
		{"Global, nonglobal middleware", "GET", "/globalafter", "two", http.StatusOK},
		{"Middleware context", "GET", "/singlectx", "first", http.StatusOK},
		{"Middleware context order", "GET", "/doublectx", "firstsecond", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Method: tt.method, URL: &url.URL{Path: tt.path}}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Body.String() != tt.body {
				t.Errorf(`%s: failed. Expected "%s", got "%s"`, tt.name, tt.body, w.Body.String())
			}

			if w.Code != tt.status {
				t.Errorf(`%s: failed. Wanted "%d", got "%d"`, tt.name, tt.status, w.Code)
			}
		})
	}
}
EOF
go test -v
rm $TESTFILE

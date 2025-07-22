package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"github.com/addyreal/simple-router"
)

func main() {
	_headers := func(w http.ResponseWriter) {
		log.Println("Headers")
	}
	_notfound := func(w http.ResponseWriter, r *http.Request) {
		log.Println("404")
	}
	_recovery := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Panic")
	}
	_global := func(n http.HandlerFunc) http.HandlerFunc {
		log.Println("Global")
		return n
	}
	_first := func(n http.HandlerFunc) http.HandlerFunc {
		log.Println("First")
		return n
	}
	_second := func(n http.HandlerFunc) http.HandlerFunc {
		log.Println("Second")
		return n
	}
	_handler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling")
	}
	_badhandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("I will panic")
		panic("")
	}


	_requests := []*http.Request{
		&http.Request{Method: "GET", URL: &url.URL{Path: "/ohno"}},
		&http.Request{Method: "GET", URL: &url.URL{Path: "/bad/ohno"}},
		&http.Request{Method: "GET", URL: &url.URL{Path: "/hello"}},
	}

	b := router.Init()
	b.SetHeaders(_headers)
	b.SetNotFound(_notfound)
	b.SetRecovery(_recovery)
	b.AddMiddleware(_global)

	b.AddGet(7, "hello", _handler)
	b.AppendMiddleware([]uint{7}, _first)
	b.AppendMiddleware([]uint{7}, _second)

	b.AddGet(0, "/bad/ohno", _badhandler)

	router := b.Get()

	for _, _req := range _requests {
		log.Println("--- Requesting", _req.URL.Path)
		router(httptest.NewRecorder(), _req)
	}
}

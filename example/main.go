package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/addyreal/simple-router"
)

func main() {
	type key string
	const arbitrarykey key = "arbitrary"
	_headers := func(w http.ResponseWriter) {
		log.Println("Headers")
	}
	_notfound := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "missing", http.StatusNotFound)
	}
	_recovery := func(err any, w http.ResponseWriter, r *http.Request) {
		switch v := err.(type) {
			case string:
				log.Println(err)
			case error:
				log.Println(v.Error())
			default:
				log.Println("becuase something bad happened")
		}
		w.Write([]byte("I recovered"))
	}
	_global := func(n http.HandlerFunc) http.HandlerFunc {
		log.Println("Global")
		return n
	}
	_first := func(n http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("First")
			val := 69
			ctx := context.WithValue(r.Context(), arbitrarykey, val)
			log.Println("I passed down", val)
			n.ServeHTTP(w, r.WithContext(ctx))
		}
	}
	_second := func(n http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Second")
			val := r.Context().Value(arbitrarykey)
			if val == nil {
				panic("error")
			}
			log.Println("I inherited", val.(int))
			n.ServeHTTP(w, r)
		}
	}
	_handler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling")
		w.Header().Set("Arbitrary", "header")
		w.Write([]byte("hello"))
	}
	_badhandler := func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(2) == 0 {
			log.Println("I will panic")
			panic("fatal error")
		}

		w.Write([]byte("I did not panic"))
	}

	_requests := []*http.Request{
		&http.Request{Method: "GET", URL: &url.URL{Path: "/ohno"}},
		&http.Request{Method: "GET", URL: &url.URL{Path: "/bad/ohno"}},
		&http.Request{Method: "GET", URL: &url.URL{Path: "/hello"}},
		&http.Request{Method: "GET", URL: &url.URL{Path: "/whoisthis"}},
		&http.Request{Method: "HEAD", URL: &url.URL{Path: "/hello"}},
	}

	b := router.Init()
	b.SetHeaders(_headers)
	b.SetNotFound(_notfound)
	b.SetRecovery(_recovery)
	b.AddMiddleware(_global)

	b.AddGet(7, "hello", _handler)
	b.AppendMiddleware([]uint8{7}, _first)
	b.AppendMiddleware([]uint8{7}, _second)

	b.AddGet(0, "/bad/ohno", _badhandler)

	router := b.Get()

	for _, _req := range _requests {
		log.Println("--- Requesting", _req.URL.Path, _req.Method)
		rec := httptest.NewRecorder()
		router(rec, _req)
		dump, _ := httputil.DumpResponse(rec.Result(), true)
		log.Println("- Response\n"+string(dump))
		log.Println("--- DONE")
	}
}

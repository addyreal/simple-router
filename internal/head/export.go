package head

import "net/http"

type HeaderWriter struct {
	http.ResponseWriter
}

func (_ *HeaderWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

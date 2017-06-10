package middleware

import (
	"net/http"
)

type headers struct {
	handler http.Handler
}

func NewJsonHeader(h http.Handler) http.Handler {
	return &headers{
		handler: h,
	}
}

func (headers *headers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	headers.handler.ServeHTTP(w, r)
}

package middleware

import (
	"log"
	"net/http"
)

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("METHOD: %s", r.Method)
		log.Printf("PATH: %s", r.URL.Path)
		log.Printf("IP ADDR: %s", r.RemoteAddr)

		h.ServeHTTP(w, r)
	}
}

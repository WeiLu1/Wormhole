package middleware

import (
	"net/http"

	"github.com/WeiLu1/wormhole/config"
)

func WithCORS(h http.HandlerFunc, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", c.Cors.AllowOrigins)
		w.Header().Set("Access-Control-Allow-Methods", c.Cors.AllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", c.Cors.AllowHeaders)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}

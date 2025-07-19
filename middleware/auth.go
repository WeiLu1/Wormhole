package middleware

import (
	"net/http"
	"strings"

	"github.com/WeiLu1/wormhole/utils"
)

func WithAuth(h http.HandlerFunc, jwtProcessor *utils.JWTProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if len(authToken) == 0 {
			http.Error(w, "No Authorization header", http.StatusBadRequest)
			return
		}

		authTokenRaw := strings.Split(authToken, " ")[1]
		_, err := jwtProcessor.VerifyToken(authTokenRaw)
		if err != nil {
			http.Error(w, "Authorization not valid", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

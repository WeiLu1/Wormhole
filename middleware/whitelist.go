package middleware

import (
	"errors"
	"log"
	"net"
	"net/http"
	"slices"
	"strings"
)

var ErrInvalidAddress = errors.New("invalid address")

func WithWhitelisting(h http.HandlerFunc, allowedAddresses []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr, err := getRealAddress(r)
		if err != nil {
			log.Print("IP address not valid")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if !slices.Contains(allowedAddresses, addr) {
			log.Print("IP address not allowed")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func getRealAddress(r *http.Request) (string, error) {
	remoteIp, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", ErrInvalidAddress
	}

	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastForwarded := addrs[len(addrs)-1]

		if ip := net.ParseIP(lastForwarded); ip != nil {
			remoteIp = ip.String()
		}
	}

	return remoteIp, nil
}

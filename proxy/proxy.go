package proxy

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/WeiLu1/wormhole/config"
)

type Proxy struct {
	Config *config.Config
}

func (p *Proxy) HandleProxy(w http.ResponseWriter, r *http.Request) {
	upstreamUrl, err := p.findTarget(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	remote, err := url.Parse(upstreamUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy := p.reverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func (p *Proxy) findTarget(path string) (string, error) {
	urlParts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	if len(urlParts) < 1 {
		return "", errors.New("insufficient path provided")
	}

	service, exists := p.Config.Services[urlParts[0]]
	if !exists {
		return "", errors.New("no valid upstream target")
	}

	remainingUrlParts := strings.TrimSuffix(strings.Join(urlParts[1:], "/"), "/")
	fullPath, _ := url.JoinPath(service.UpstreamPath, remainingUrlParts)

	return fullPath, nil
}

func (p *Proxy) reverseProxy(url *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		targetPath := url.Path
		req.URL.Path = targetPath
		req.URL.RawPath = targetPath
		req.Host = url.Host
	}

	return proxy
}

package server

import (
	"net/http"

	"github.com/WeiLu1/wormhole/config"
	mw "github.com/WeiLu1/wormhole/middleware"
	"github.com/WeiLu1/wormhole/proxy"
)

type Server struct {
	Config *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{
		Config: c,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	proxy := proxy.Proxy{Config: s.Config}

	mux.HandleFunc("/*", mw.WithLogging(mw.WithCORS(proxy.HandleProxy, s.Config)))

	return http.ListenAndServe(":"+s.Config.Port, mux)
}

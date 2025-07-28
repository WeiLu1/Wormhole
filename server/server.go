package server

import (
	"log"
	"net/http"
	"os"

	"github.com/WeiLu1/wormhole/config"
	mw "github.com/WeiLu1/wormhole/middleware"
	"github.com/WeiLu1/wormhole/proxy"
	"github.com/WeiLu1/wormhole/utils"
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

	rateLimitConfig := s.Config.RateLimiting
	rateLimiter := mw.NewRateLimiter()
	if rateLimitConfig.PerIpAddress.Enabled {
		go rateLimiter.CleanUpVisitors(rateLimitConfig.PerIpAddress.CleanupIntervalSeconds)
	}

	handler := proxy.HandleProxy
	handler = mw.WithCORS(handler, s.Config)

	if rateLimitConfig.PerIpAddress.Enabled {
		log.Print("Running with visitor rate limit")
		handler = mw.WithRateLimitingPerVisitor(handler, rateLimitConfig, rateLimiter)
	} else {
		log.Print("Running with global rate limit")
		handler = mw.WithRateLimitingGlobal(handler, rateLimitConfig, rateLimiter)
	}

	if s.Config.Auth.UseJwt {
		jwtProcessor := utils.NewJWTProcessor(os.Getenv("JWT_SECRET"))
		handler = mw.WithAuth(handler, jwtProcessor)
	}

	if allowAddrs := s.Config.Whitelist.Allow; allowAddrs != nil {
		handler = mw.WithWhitelisting(handler, allowAddrs)
	}

	handler = mw.WithLogging(handler)

	mux.HandleFunc("/*", handler)

	log.Print("Running server...")
	return http.ListenAndServe(":"+s.Config.Port, mux)
}

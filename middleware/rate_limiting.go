package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/WeiLu1/wormhole/config"
)

type Visitor struct {
	lastSeen    time.Time
	leakyBucket *LeakyBucket
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*Visitor),
	}
}

func (rl *RateLimiter) CleanUpVisitors(cleanupIntervalSeconds int) {
	ticker := time.NewTicker(time.Duration(cleanupIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > time.Duration(cleanupIntervalSeconds)*time.Second {
				log.Printf("Removing client %s", ip)
				v.leakyBucket.stop <- struct{}{}
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getOrCreateVisitorBucket(ip string, maxCapacity int64) *LeakyBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			lastSeen:    time.Now(),
			leakyBucket: NewLeakyBucket(maxCapacity),
		}
		rl.visitors[ip] = visitor
	} else {
		visitor.lastSeen = time.Now()
	}

	return visitor.leakyBucket
}

type LeakyBucket struct {
	maxCapacity int64
	queue       chan struct{}
	stop        chan struct{}
}

func NewLeakyBucket(maxCapacity int64) *LeakyBucket {
	lb := &LeakyBucket{
		maxCapacity: maxCapacity,
		queue:       make(chan struct{}, maxCapacity),
		stop:        make(chan struct{}),
	}

	go lb.leak()
	return lb
}

// Checks whether or not a request can be processed by seeing if the capacity of the queue channel is full.
func (lb *LeakyBucket) Allow() bool {
	select {
	case lb.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

// This will remove a token from the queue in a time interval specified by the maximum amount of requests provided in the configuration per second.
// Leaking will stop when the stop channel recieves
func (lb *LeakyBucket) leak() {
	ticker := time.NewTicker(time.Second / time.Duration(lb.maxCapacity))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case <-lb.queue:
			default:
			}
		case <-lb.stop:
			return
		}
	}
}

func WithRateLimitingGlobal(h http.HandlerFunc, c config.RateLimiting, rateLimiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serve(rateLimiter, "global", c, w, h, r)
	}
}

func WithRateLimitingPerVisitor(h http.HandlerFunc, c config.RateLimiting, rateLimiter *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		serve(rateLimiter, ip, c, w, h, r)
	}
}

func serve(rateLimiter *RateLimiter, ip string, c config.RateLimiting, w http.ResponseWriter, h http.HandlerFunc, r *http.Request) {
	lb := rateLimiter.getOrCreateVisitorBucket(ip, c.MaxCapacity)
	if !lb.Allow() {
		rateLimitExceededResponse(w)
		return
	}

	h.ServeHTTP(w, r)
}

func rateLimitExceededResponse(w http.ResponseWriter) {
	log.Print("Too many requests for rate limiting capacity")
	http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
}

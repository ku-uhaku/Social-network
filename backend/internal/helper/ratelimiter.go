package helper

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	window   time.Duration
	limits   map[string]int
	requests map[string][]time.Time
	mu       sync.Mutex
}

func Neewratelimeter(window time.Duration) *RateLimiter {
	return &RateLimiter{
		window: window,
		limits: map[string]int{
			"api":      300,
			"authonti": 20,
		},
		requests: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Wraponall(limitType string, next http.Handler) http.Handler {
	limit, ok := rl.limits[limitType]
	if !ok {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Invalid rate limit type", http.StatusInternalServerError)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions && !rl.allow(clientIP(r), limit) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string, limit int) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	requests := rl.requests[ip]

	// Remove requests outside the current window.
	cutoff := now.Add(-rl.window)
	validRequests := requests[:0]

	for _, requestTime := range requests {
		if requestTime.After(cutoff) {
			validRequests = append(validRequests, requestTime)
		}
	}

	if len(validRequests) >= limit {
		rl.requests[ip] = validRequests
		return false
	}

	rl.requests[ip] = append(validRequests, now)

	return true
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
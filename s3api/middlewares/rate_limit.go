package middlewares

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter holds per-IP request counters with a sliding window.
type RateLimiter struct {
	mu       sync.Mutex
	counters map[string]*rateBucket
	limit    int
	window   time.Duration
}

type rateBucket struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a RateLimiter that allows up to limit requests
// per window duration per remote address.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		counters: make(map[string]*rateBucket),
		limit:    limit,
		window:   window,
	}
}

// Allow returns true if the request from addr is within the rate limit.
func (rl *RateLimiter) Allow(addr string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, ok := rl.counters[addr]
	if !ok || now.After(bucket.windowEnd) {
		rl.counters[addr] = &rateBucket{
			count:     1,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	if bucket.count >= rl.limit {
		return false
	}
	bucket.count++
	return true
}

// RateLimitMiddleware returns an HTTP middleware that enforces the rate limit.
func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addr := r.RemoteAddr
			if !rl.Allow(addr) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

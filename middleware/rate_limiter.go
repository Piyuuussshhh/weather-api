package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rateLimiterEntry
	mu sync.Mutex
	rate rate.Limit
	burst int
}

type rateLimiterEntry struct {
    limiter    *rate.Limiter
    lastActive time.Time
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rateLimiterEntry),
		rate:     r,
		burst:    burst,
	}

	// Background cleanup goroutine.
	go rl.cleanupExpiredLimiters()

	return rl
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		// Use the remote address as the key
		ip := r.RemoteAddr 
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}	

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if entry, exists := rl.limiters[ip]; exists {
		entry.lastActive = time.Now()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[ip] = &rateLimiterEntry{
		limiter:    limiter,
		lastActive: time.Now(),
	}

	return limiter
}

func (rl *RateLimiter) cleanupExpiredLimiters() {
    for {
		// Run cleanup every minute
        time.Sleep(1 * time.Minute) 

        rl.mu.Lock()
        for ip, entry := range rl.limiters {
            if time.Since(entry.lastActive) > 10*time.Minute {
                delete(rl.limiters, ip)
            }
        }
        rl.mu.Unlock()
    }
}
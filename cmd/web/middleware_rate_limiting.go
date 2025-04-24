package main

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// Define a struct to hold rate limiters for each IP
type ipRateLimiter struct {
	mu sync.Mutex
	limiters map[string]*rate.Limiter
}

func newIPRateLimiter(rateLimit rate.Limit, burst int) *ipRateLimiter {
	return &ipRateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (irl *ipRateLimiter) getLimiter(ip string, rateLimit rate.Limit, burst int) *rate.Limiter {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	limiter, exists := irl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burst)
		irl.limiters[ip] = limiter
	}
	return limiter
}

func (app *application) rateLimitMiddleware(next http.Handler) http.Handler {
	// Define your rate limit (e.g., 10 requests per second with a burst of 20)
	rateLimit := rate.Limit(10)
	burst := 20
	ipLimiters := newIPRateLimiter(rateLimit, burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Or a more sophisticated way to identify users

		limiter := ipLimiters.getLimiter(ip, rateLimit, burst)

		if !limiter.Allow() {
			app.rateLimitExceeded(w,r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "1") // Suggest the client to retry after 1 second
	data := NewTemplateData()
	data.AlertMessage = "Too many requests. Please try again later."
	data.AlertType = "alert-danger"
	app.render(w,r, http.StatusTooManyRequests, "error-404.tmpl", data) // Reuse or create a specific error template
}
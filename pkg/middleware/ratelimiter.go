package reflectormiddleware

import (
	"net/http"
	"sync"
	"time"
)

func RateLimiter(interval time.Duration, limit int) func(http.Handler) http.Handler {
	type client struct {
		count    int
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{}
			}
			c := clients[ip]
			now := time.Now()
			if now.Sub(c.lastSeen) > interval {
				c.count = 0
				c.lastSeen = now
			}
			c.count++
			if c.count > limit {
				mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

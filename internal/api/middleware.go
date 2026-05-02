package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// We use a map to track limiters for every unique IP address
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// RateLimiter prevents spamming the "Pulse" endpoints
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := visitors[ip]
		if !exists {
			// Allow 1 request per second (r), with a burst of 5 (b)
			limiter = rate.NewLimiter(1, 5)
			visitors[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Pulse limit exceeded. Slow down!",
			})
			return
		}
		c.Next()
	}
}

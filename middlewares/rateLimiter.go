// middlewares/rateLimiter.go
package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiters = make(map[string]*rateLimiter)
	mu           sync.Mutex
)

type rateLimiter struct {
	tokens        int
	lastTokenTime time.Time
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		tokens:        5,
		lastTokenTime: time.Now(),
	}
}

func RateLimitMiddleware(c *gin.Context) {
	mu.Lock()
	limiter, exists := rateLimiters[c.ClientIP()]
	if !exists {
		limiter = newRateLimiter()
		rateLimiters[c.ClientIP()] = limiter
	}
	mu.Unlock()

	limiter.tokens += int(time.Since(limiter.lastTokenTime).Seconds())
	if limiter.tokens > 5 {
		limiter.tokens = 5
	}
	limiter.lastTokenTime = time.Now()

	if limiter.tokens > 0 {
		limiter.tokens--
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
	}
}

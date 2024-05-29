package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ryuuukin/ap-assignment1/middlewares"
	"github.com/gin-gonic/gin"
)

func TestRateLimitMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middlewares.RateLimitMiddleware)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send multiple requests to test rate limiting
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			t.Logf("Request %d: Rate limit exceeded", i+1)
		} else if w.Code == http.StatusOK {
			t.Logf("Request %d: Success", i+1)
		} else {
			t.Errorf("Request %d: Unexpected status code: %d", i+1, w.Code)
		}
	}

}

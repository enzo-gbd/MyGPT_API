package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestLimiterMiddleware(t *testing.T) {
	r := rate.Every(time.Hour)
	limiter := rate.NewLimiter(r, 1)

	gin.SetMode(gin.TestMode)
	setupRouter()
	router.Use(Limiter(limiter))
	defer sqlDB.Close()

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d for rate limit, got %d", http.StatusTooManyRequests, w.Code)
	}
}

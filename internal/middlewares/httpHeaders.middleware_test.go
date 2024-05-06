package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPHeadersMiddleware(t *testing.T) {
	setupRouter()
	router.Use(HTTPHeaders())
	defer sqlDB.Close()

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, "default-src 'self'", w.Header().Get("Content-Security-Policy"), "Content-Security-Policy should be correctly set")
	assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"), "Strict-Transport-Security should be correctly set")
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"), "X-Content-Type-Options should be correctly set")
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"), "X-Frame-Options should be correctly set")
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"), "Cache-Control should be correctly set")
	assert.Equal(t, "GBA", w.Header().Get("Server"), "Server should be correctly set")
}

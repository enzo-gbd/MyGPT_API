// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import (
	"net/http"

	"github.com/enzo-gbd/GBA/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Limiter returns a middleware handler function that limits the rate of incoming requests.
// It takes a *rate.Limiter object, which defines the maximum rate at which requests
// can be processed.
//
// If a request arrives and the rate limit has been exceeded, the middleware
// aborts the request with an HTTP 429 (Too Many Requests) status. If the rate
// limit has not been exceeded, the middleware allows the request to proceed.
func Limiter(limiter *rate.Limiter) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !limiter.Allow() {
			utils.AbortWithError(context, http.StatusTooManyRequests, "Too many requests")
			return
		}
		context.Next()
	}
}

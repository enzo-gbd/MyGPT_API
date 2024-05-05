// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors returns a middleware handler that implements CORS (Cross-Origin Resource Sharing) controls.
// This middleware allows for configuring various aspects of CORS such as allowed origins,
// methods, headers, and credentials.
//
// The configuration allows requests from "https://localhost" and supports methods
// GET, POST, PUT, and DELETE. It also handles specific headers and exposes some,
// allowing credentials and supporting the use of wildcards in origin specification.
// The `MaxAge` setting controls how long the results of a preflight request can be cached.
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * time.Hour,
	})
}

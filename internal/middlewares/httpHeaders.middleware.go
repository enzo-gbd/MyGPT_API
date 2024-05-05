// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import "github.com/gin-gonic/gin"

// HTTPHeaders returns a middleware handler function that sets several security-related
// HTTP headers on responses. This middleware is intended to enhance the security
// of the application by specifying browser behavior regarding content sources,
// transport security, frame options, content sniffing, and caching.
//
// The set headers are:
//   - Content-Security-Policy: Restricts resources the client can load for a given page,
//     with 'default-src 'self” meaning that all resources must come from the same origin.
//   - Strict-Transport-Security: Enforces secure (HTTPS) connections to the server and
//     includes subdomains for a maximum of one year.
//   - X-Content-Type-Options: Prevents the browser from interpreting files as a
//     different MIME type than what is specified by the content type in the HTTP headers.
//   - X-Frame-Options: Prevents clickjacking attacks by preventing content from being
//     embedded into other sites.
//   - Cache-Control: Directs browsers not to store the page in cache.
//   - Server: Customizes the server name sent in the HTTP response header.
func HTTPHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Server", "GBA")
		c.Next()
	}
}

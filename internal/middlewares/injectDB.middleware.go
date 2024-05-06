// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InjectDB creates a middleware that injects a *gorm.DB instance into the Gin context.
// This allows for easy access to the database connection throughout the request's lifecycle
// in subsequent handlers.
//
// The `db` parameter is the GORM database object that should be accessible
// in the handlers after this middleware is applied.
func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

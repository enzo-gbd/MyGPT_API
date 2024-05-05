// Package utils provides utility functions that support various operations across the application.
package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDatabaseInContext retrieves the database instance from the given Gin context.
// It attempts to extract a *gorm.DB object stored with the key "db". If the database instance
// is not found in the context, an error is returned.
//
// Parameters:
//   - context: *gin.Context representing the current request context.
//
// Returns:
//   - *gorm.DB: A pointer to the GORM database object if found.
//   - error: An error object that reports an absence of the database in the context.
func GetDatabaseInContext(context *gin.Context) (*gorm.DB, error) {
	database, exists := context.Get("db")
	if !exists {
		return nil, errors.New("database not available")
	}
	return database.(*gorm.DB), nil
}

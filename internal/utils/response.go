// Package utils provides utility functions that support various operations across the application.
package utils

import (
	"github.com/gin-gonic/gin"
)

// AbortWithError sends a JSON response with a failure status and a custom message.
// It uses the specified HTTP status code and includes the provided message in the response.
// This function also aborts the request chain, meaning no subsequent handlers will be executed.
//
// Parameters:
//
//	c       - The Gin context to use for sending the response.
//	code    - The HTTP status code to send.
//	message - The error message to include in the JSON response.
func AbortWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": "fail", "message": message})
}

// SendSuccess sends a JSON response containing the provided data.
// This function uses the specified HTTP status code for the response.
//
// Parameters:
//
//	c    - The Gin context to use for sending the response.
//	code - The HTTP status code to use for the response.
//	data - The data to include in the JSON response.
func SendSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

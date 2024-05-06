// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import (
	"net/http"

	"github.com/enzo-gbd/GBA/internal/utils"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/gin-gonic/gin"
)

// CheckUserRole returns a middleware handler function that checks if the
// currently logged-in user has the required role to proceed with the request.
// It takes a string parameter `role`, which specifies the required user role.
//
// This middleware extracts the current user from the context and validates
// their role. If the user is not logged in, or if the user data is not of type *models.User,
// or if the user's role does not match the required role, it aborts the request
// with an HTTP status of 401 (Unauthorized).
func CheckUserRole(role string) gin.HandlerFunc {
	return func(context *gin.Context) {
		value, exists := context.Get("currentUser")
		if !exists {
			utils.AbortWithError(context, http.StatusUnauthorized, "You are not logged in")
			return
		}

		currentUser, ok := value.(*models.User)
		if !ok {
			utils.AbortWithError(context, http.StatusUnauthorized, "invalid user type")
			return
		}

		if currentUser.Role != role {
			utils.AbortWithError(context, http.StatusUnauthorized, "You are not allowed")
			return
		}
	}
}

// Package middlewares contains middleware functions for handling various
// aspects of HTTP requests within the application.
package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/gin-gonic/gin"
)

// DeserializeUser returns a middleware handler function that processes the
// authentication token (JWT) provided in the request's Authorization header
// or in a cookie named 'access_token'. It validates the token, retrieves the
// user associated with the token from the database, and sets this user in the
// request context for use in subsequent request handling.
//
// The function first attempts to extract the token from the 'Authorization'
// header. If not found, it then checks for the token in a cookie. If neither
// are present or valid, it aborts the request with an HTTP status of 401
// (Unauthorized).
//
// On successful token validation and user retrieval, the user is set in the
// request context under the key 'currentUser'. If any step fails, the function
// aborts the process, providing appropriate HTTP error responses, including
// 500 (Internal Server Error) for server-related issues and 404 (Not Found)
// if the user does not exist in the database.
func DeserializeUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		database, err := utils.GetDatabaseInContext(context)
		if err != nil {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
			return
		}
		var accessToken string
		cookie, _ := context.Cookie("access_token")

		authorizationHeader := context.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			utils.AbortWithError(context, http.StatusUnauthorized, "You are not logged in")
			return
		}

		config, _ := configs.LoadConfig()
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			utils.AbortWithError(context, http.StatusUnauthorized, "The access token is not valid")
			return
		}

		var user *models.User
		result := database.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			utils.AbortWithError(context, http.StatusNotFound, "the user belonging to this token no longer exists")
			return
		}

		context.Set("currentUser", user)
		context.Next()
	}
}

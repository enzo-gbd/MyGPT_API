// Package api provides the routing functionalities
// It sets up routes and associates them with their respective handlers.
package api

import (
	"github.com/enzo-gbd/GBA/internal/controllers/auth"
	"github.com/enzo-gbd/GBA/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// AuthRouteController holds a reference to an auth.AuthController to handle authentication-related requests.
type AuthRouteController struct {
	authController auth.AuthController
}

// NewAuthRouteController creates a new instance of AuthRouteController with the provided authController.
// It returns an AuthRouteController which can be used to set up authentication routes.
func NewAuthRouteController(authController auth.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

// AuthRoutes sets up the routing for all authentication-related endpoints under the provided RouterGroup.
// It registers routes for user registration, login, token refresh, and logout.
func (ac *AuthRouteController) AuthRoutes(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/register", ac.authController.SignUpUser)                              // Registers a new user.
	router.POST("/login", ac.authController.SignInUser)                                 // Authenticates a user and returns a session token.
	router.POST("/refresh", ac.authController.RefreshAccessToken)                       // Refreshes an existing session token.
	router.POST("/logout", middlewares.DeserializeUser(), ac.authController.LogoutUser) // Ends a user's session.
}

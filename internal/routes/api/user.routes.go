// Package api provides the routing functionalities
// It sets up routes and associates them with their respective handlers.
package api

import (
	"github.com/enzo-gbd/GBA/internal/controllers/user"
	"github.com/enzo-gbd/GBA/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// UserAPIRouteController handles the routing of user-related API endpoints.
type UserAPIRouteController struct {
	userController user.UserController
}

// NewAPIRouteUserController creates a new instance of UserAPIRouteController
// using the provided userController. This controller will be used to handle
// user-specific requests.
func NewAPIRouteUserController(userController user.UserController) UserAPIRouteController {
	return UserAPIRouteController{userController}
}

// UserRoute configures the routes related to the user in the provided RouterGroup.
// It sets up middleware for deserializing the user and defines the "me" route to
// fetch the user's own data.
func (uc *UserAPIRouteController) UserRoute(rg *gin.RouterGroup) {
	rg.Use(middlewares.DeserializeUser())   // Apply middleware to deserialize user information from incoming requests.
	router := rg.Group("me")                // Group routes under 'me' for current user operations.
	router.GET("", uc.userController.GetMe) // Define the GET request for 'me' to fetch current user's data.
}

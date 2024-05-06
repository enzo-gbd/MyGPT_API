// Package admin provides route controllers for managing operations in an administrative context.
package admin

import (
	"github.com/enzo-gbd/GBA/internal/controllers/user"
	"github.com/gin-gonic/gin"
)

// UserAdminRouteController handles the routing of user administration functions.
type UserAdminRouteController struct {
	userController user.UserController // userController manages the basic user operations.
}

// NewAdminRouteUserController creates a new instance of UserAdminRouteController using the provided userController.
// This setup allows for dependency injection and easier management of user-related routes.
func NewAdminRouteUserController(userController user.UserController) UserAdminRouteController {
	return UserAdminRouteController{userController}
}

// UserRoute defines routes for user management within an admin-specific router group.
// The paths include operations to retrieve all users, retrieve a user by ID, update a user, and delete a user.
func (uc *UserAdminRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("/", uc.userController.GetUsers)         // GetUsers handles the retrieval of all users.
	router.GET("/:id", uc.userController.GetUserByID)   // GetUserByID handles fetching a specific user based on user ID.
	router.PUT("/:id", uc.userController.UpdateUser)    // UpdateUser handles updating a specific user's details.
	router.DELETE("/:id", uc.userController.DeleteUser) // DeleteUser handles the removal of a user by ID.
}

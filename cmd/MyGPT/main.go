// Package main provides the entry point for the GBA server.
// It sets up routing, middleware, and starts the server with TLS support.
package main

import (
	"log"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/controllers/auth"
	"github.com/enzo-gbd/GBA/internal/controllers/user"
	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/enzo-gbd/GBA/internal/middlewares"
	"github.com/enzo-gbd/GBA/internal/routes/admin"
	"github.com/enzo-gbd/GBA/internal/routes/api"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// AuthRouteController handles the authentication-related routes.
	AuthRouteController api.AuthRouteController

	// UserAPIRouteController handles user management within the API scope.
	UserAPIRouteController api.UserAPIRouteController

	// UserAdminRouteController handles user management within the admin scope.
	UserAdminRouteController admin.UserAdminRouteController
)

// init initializes the controllers for the API and administration routes.
func init() {
	authController := auth.NewAuthController()
	AuthRouteController = api.NewAuthRouteController(authController)

	userController := user.NewUserController()
	UserAPIRouteController = api.NewAPIRouteUserController(userController)
	UserAdminRouteController = admin.NewAdminRouteUserController(userController)
}

// apiRoutes configures the API and admin routes with the appropriate controllers and middleware.
func apiRoutes(router *gin.Engine) {
	apiRouter := router.Group("/api")
	{
		AuthRouteController.AuthRoutes(apiRouter)
		UserAPIRouteController.UserRoute(apiRouter)
	}
	adminRouter := router.Group("/admin")
	adminRouter.Use(middlewares.DeserializeUser())
	adminRouter.Use(middlewares.CheckUserRole("admin"))
	{
		UserAdminRouteController.UserRoute(adminRouter)
	}
}

// setupRouter initializes the Gin engine with middleware, database, and rate limiting.
// It returns the configured router.
func setupRouter(config *configs.Config) *gin.Engine {
	router := gin.Default()
	database := db.InitDB(config)

	limiter := rate.NewLimiter(1, 5)

	router.Use(middlewares.Cors())
	router.Use(middlewares.HTTPHeaders())
	router.Use(middlewares.Limiter(limiter))
	router.Use(middlewares.InjectDB(database))

	return router
}

// main is the entry point of the application.
// It loads the configuration, sets up the router, and starts the server.
func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Could not load environment variables: ", err)
	}
	router := setupRouter(&config)

	apiRoutes(router)

	err = router.RunTLS(":8443", "tls/cert.pem", "tls/key.pem")
	if err != nil {
		log.Fatal("Failed to start HTTPS server: ", err)
	}
}

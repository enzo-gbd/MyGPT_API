package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() AuthController {
	return AuthController{}
}

// SignUpUser handles user registration.
// @Summary Register a new user
// @Description Registers a new user with the necessary details provided in the request body.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body models.SignUpInput true "User Registration Data"
// @Success 201 {object} map[string]interface{} "Returns status success on successful registration"
// @Failure 400 {object} map[string]interface{} "Returns error message for bad request"
// @Failure 409 {object} map[string]interface{} "Returns error message for email already exists"
// @Failure 500 {object} map[string]interface{} "Returns error message for internal server error"
// @Router /register [post]
func (ac *AuthController) SignUpUser(context *gin.Context) {
	var payload *models.SignUpInput
	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	if err := context.ShouldBindJSON(&payload); err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}

	err = payload.Validate()
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	newUser := models.User{
		FirstName: payload.FirstName,
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Birthday:  payload.Birthday,
		Gender:    payload.Gender,
		Role:      "user",
		Verified:  false,
	}
	err = newUser.Validate()
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}

	var exists models.User
	if err := database.Where("email = ?", newUser.Email).First(&exists).Error; err == nil {
		utils.AbortWithError(context, http.StatusConflict, "User with that email already exists")
		return
	}

	result := database.Create(&newUser)

	if result.Error != nil {
		utils.AbortWithError(context, http.StatusBadGateway, "Something bad happened")
		return
	}

	utils.SendSuccess(context, http.StatusCreated, gin.H{"status": "success"})
}

// SignInUser handles user login.
// @Summary Login a user
// @Description Logs in a user and returns an access token and a refresh token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body models.SignInInput true "User Login Data"
// @Success 200 {object} map[string]interface{} "Returns the access token and refresh token on successful login"
// @Failure 400 {object} map[string]interface{} "Returns error message for bad request"
// @Failure 401 {object} map[string]interface{} "Returns error message for invalid email or password"
// @Failure 500 {object} map[string]interface{} "Returns error message for internal server error"
// @Router /login [post]
func (ac *AuthController) SignInUser(context *gin.Context) {
	var payload *models.SignInInput
	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	if err := context.ShouldBindJSON(&payload); err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}

	err = payload.Validate()
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	result := database.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		utils.AbortWithError(context, http.StatusUnauthorized, "Invalid email or Password")
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		utils.AbortWithError(context, http.StatusUnauthorized, "Invalid email or Password")
		return
	}

	config, _ := configs.LoadConfig()
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, "Configuration error")
		return
	}

	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := utils.GenerateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	context.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	context.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	context.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.SendSuccess(context, http.StatusOK, gin.H{"status": "success", "token": accessToken})
}

// LogoutUser handles user logout.
// @Summary Logout a user
// @Description Logs out a user by invalidating the tokens.
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns status success on successful logout"
// @Router /logout [post]
func (ac *AuthController) LogoutUser(context *gin.Context) {
	context.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	context.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	context.SetCookie("logged_in", "", -1, "/", "localhost", false, false)
	utils.SendSuccess(context, http.StatusOK, gin.H{"status": "success"})
}

// RefreshAccessToken refreshes the access token using a refresh token.
// @Summary Refresh access token
// @Description Refreshes an access token using the refresh token provided in cookie.
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns new access token"
// @Failure 401 {object} map[string]interface{} "Returns error message for unauthorized or invalid token"
// @Failure 404 {object} map[string]interface{} "Returns error message when the user does not exist"
// @Failure 500 {object} map[string]interface{} "Returns error message for internal server error"
// @Router /refresh [post]
func (ac *AuthController) RefreshAccessToken(context *gin.Context) {
	cookie, err := context.Cookie("refresh_token")
	if err != nil {
		utils.AbortWithError(context, http.StatusUnauthorized, "could not found refresh token")
		return
	}

	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	config, _ := configs.LoadConfig()
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, "Configuration error")
		return
	}

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		utils.AbortWithError(context, http.StatusUnauthorized, "The refresh token is not valid")
		return
	}

	var user models.User
	result := database.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		utils.AbortWithError(context, http.StatusNotFound, "the user belonging to this token no logger exists")
		return
	}

	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}

	context.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	context.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	utils.SendSuccess(context, http.StatusOK, gin.H{"status": "success", "token": accessToken})
}

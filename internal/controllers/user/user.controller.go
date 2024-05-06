package user

import (
	"errors"
	"net/http"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

// GetUsers retrieves all users.
// @Summary Get all users
// @Description Fetches all users from the database.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /users [get]
func (uc *UserController) GetUsers(context *gin.Context) {
	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	var users []models.User

	if err := database.Find(&users).Error; err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	if len(users) == 0 {
		utils.AbortWithError(context, http.StatusNotFound, "can't found any users")
		return
	}
	utils.SendSuccess(context, http.StatusOK, users)
}

// GetUserByID fetches a single user by UUID.
// @Summary Get user by ID
// @Description Fetches a user from the database by UUID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(context *gin.Context) {
	idStr := context.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	var user models.User

	if err := database.Where("id = ?", idStr).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found user")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.SendSuccess(context, http.StatusOK, user)
}

// UpdateUser updates an existing user.
// @Summary Update user
// @Description Updates details of an existing user.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User Data"
// @Success 200 {object} models.User
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(context *gin.Context) {
	idStr := context.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}
	err = user.Validate()
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}
	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	var oldUser models.User
	if err := database.Where("id = ?", idStr).First(&oldUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found user")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	user.ID = oldUser.ID
	if err := database.Save(&user).Error; err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(context, http.StatusOK, user)
}

// DeleteUser deletes a user.
// @Summary Delete user
// @Description Deletes a user by UUID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(context *gin.Context) {
	idStr := context.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	var user models.User

	if err := database.Where("id = ?", idStr).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found user")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := database.Delete(&user).Error; err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(context, http.StatusOK, gin.H{})
}

// GetMe retrieves the current logged in user's information.
// @Summary Get current user
// @Description Retrieves information about the current logged in user.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} object
// @Router /users/me [get]
func (uc *UserController) GetMe(context *gin.Context) {
	obj, exists := context.Get("currentUser")
	if !exists {
		utils.AbortWithError(context, http.StatusUnauthorized, "You are not logged in")
		return
	}

	currentUser := obj.(*models.User)
	userResponse := &models.UserResponse{
		ID:               currentUser.ID,
		Name:             currentUser.Name,
		Birthday:         currentUser.Birthday,
		Gender:           currentUser.Gender,
		Email:            currentUser.Email,
		Role:             currentUser.Role,
		Address:          currentUser.Address.String,
		SubscriptionCode: currentUser.SubscriptionCode.String,
		CreatedAt:        currentUser.CreatedAt,
		UpdatedAt:        currentUser.UpdatedAt,
	}

	utils.SendSuccess(context, http.StatusOK, userResponse)
}

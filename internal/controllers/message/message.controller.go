package message

import (
	"errors"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

type MessageController struct{}

func NewMessageController() MessageController { return MessageController{} }

// GetMessageByID fetches a single message by UUID.
// @Summary Get message by ID
// @Description Fetches a message from the database by UUID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /messages/{id} [get]
func (mc MessageController) GetMessageByID(context *gin.Context) {
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
	var message models.Message

	if err := database.Where("id = ?", idStr).First(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found message")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.SendSuccess(context, http.StatusOK, message)
}

// UpdateMessage updates an existing user.
// @Summary Update message
// @Description Updates details of an existing message.
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param user body models.Message true "Message Data"
// @Success 200 {object} models.Message
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /messages/{id} [put]
func (mc *MessageController) UpdateMessage(context *gin.Context) {
	idStr := context.Param("id")

	_, err := uuid.Parse(idStr)
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	var message models.Message
	if err := context.ShouldBindJSON(&message); err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}
	err = message.Validate()
	if err != nil {
		utils.AbortWithError(context, http.StatusBadRequest, err.Error())
		return
	}
	database, err := utils.GetDatabaseInContext(context)
	if err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	var oldMessage models.Message
	if err := database.Where("id = ?", idStr).First(&oldMessage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found message")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}

	message.ID = oldMessage.ID
	if err := database.Save(&message).Error; err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(context, http.StatusOK, message)
}

// DeleteMessage deletes a message.
// @Summary Delete message
// @Description Deletes a message by UUID.
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /messages/{id} [delete]
func (mc *MessageController) DeleteMessage(context *gin.Context) {
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
	var message models.Message

	if err := database.Where("id = ?", idStr).First(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.AbortWithError(context, http.StatusNotFound, "Can't found message")
		} else {
			utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := database.Delete(&message).Error; err != nil {
		utils.AbortWithError(context, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendSuccess(context, http.StatusOK, gin.H{})
}

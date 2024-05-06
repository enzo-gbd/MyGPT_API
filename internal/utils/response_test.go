package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	AbortWithError(c, http.StatusBadRequest, "error message")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRespondWithSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendSuccess(c, http.StatusOK, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

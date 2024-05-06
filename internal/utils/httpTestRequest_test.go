package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/test", func(context *gin.Context) {
		var payload map[string]interface{}
		if err := context.ShouldBindJSON(&payload); err != nil {
			AbortWithError(context, http.StatusBadRequest, err.Error())
			return
		}
		SendSuccess(context, http.StatusOK, payload)
	})
	return router
}

func TestStructToIOReader(t *testing.T) {
	type testStruct struct {
		Test string `json:"test"`
	}
	testObj := testStruct{Test: "value"}

	reader, err := StructToIOReader(testObj)
	assert.NoError(t, err)

	// Lire les données du io.Reader pour vérifier la conversion
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return
	}
	assert.Equal(t, `{"test":"value"}`, buf.String())
}

func TestHttpTestRequest(t *testing.T) {
	router := SetupTestRouter()
	testObj := map[string]interface{}{"test": "value"}

	w, err := HttpTestRequest(router, "POST", "/test", testObj)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifier le corps de la réponse
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "value", response["test"])
}

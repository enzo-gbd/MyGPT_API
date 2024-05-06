package utils

import (
	"errors"
	"testing"

	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseInContext(t *testing.T) {
	database, _, _ := db.InitMockDB()
	context := gin.Context{}

	context.Set("db", database)

	retrievedDb, err := GetDatabaseInContext(&context)

	assert.Nil(t, err)
	assert.Equal(t, database, retrievedDb, "the database in the context is not same as the initial database")

	context = gin.Context{}

	retrievedDb, err = GetDatabaseInContext(&context)

	assert.NotNil(t, err)
	assert.Nil(t, retrievedDb)
	assert.Equal(t, err, errors.New("database not available"), "no error when the database is not retrieved")
}

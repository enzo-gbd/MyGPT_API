package middlewares

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine
var database *gorm.DB
var sqlDB *sql.DB
var mock sqlmock.Sqlmock

func setupRouter() {
	router = gin.Default()
	database, sqlDB, mock = db.InitMockDB()

	router.Use(InjectDB(database))
}

func TestMain(m *testing.M) {
	m.Run()
}

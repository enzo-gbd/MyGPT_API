package middlewares

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestInjectDBMiddleware(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db, %s", err)
	}

	c, _ := gin.CreateTestContext(nil)

	InjectDB(db)(c)

	value, exists := c.Get("db")
	assert.True(t, exists, "Le contexte doit contenir l'objet DB")
	assert.NotNil(t, value, "L'objet DB ne doit pas être nil")
	assert.IsType(t, &gorm.DB{}, value, "L'objet doit être de type *gorm.DB")
}

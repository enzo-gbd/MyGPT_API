package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/models/builders"
	"github.com/gin-gonic/gin"
)

func TestCheckUserRole(t *testing.T) {
	method, url := "GET", "/"
	tests := []struct {
		name         string
		initRouter   func(*gin.Engine)
		expectedCode int
	}{
		{
			name: "With admin user",
			initRouter: func(r *gin.Engine) {
				r.GET(url, setCurrentUser(
					builders.NewUserBuilder().WhereRole("admin").Build(),
				), CheckUserRole("admin"), func(c *gin.Context) {})
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "With not admin user",
			initRouter: func(r *gin.Engine) {
				r.GET(url, setCurrentUser(
					builders.NewUserBuilder().Build(),
				), CheckUserRole("admin"), func(c *gin.Context) {})
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "With no user",
			initRouter: func(r *gin.Engine) {
				r.GET(url, setCurrentUser(models.User{}), CheckUserRole("admin"), func(c *gin.Context) {})
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "With no 'currentUser' set",
			initRouter: func(r *gin.Engine) {
				r.GET(url, CheckUserRole("admin"), func(c *gin.Context) {})
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRouter()
			tt.initRouter(router)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, url, nil)
			router.ServeHTTP(w, req)
			defer sqlDB.Close()

			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}
		})
	}
}

func setCurrentUser(user models.User) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("currentUser", &user)
	}
}

package middlewares

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/models/builders"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/enzo-gbd/GBA/internal/utils/testUtils"
	"github.com/gin-gonic/gin"
)

func TestDeserializeUser(t *testing.T) {
	queryFirst := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`
	tests := []struct {
		name          string
		tokenInCookie bool
		tokenInHeader bool
		invalidToken  bool
		expectedCode  int
	}{
		{
			name:          "Valid token in cookie and header",
			tokenInCookie: true,
			tokenInHeader: true,
			invalidToken:  false,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Valid token in cookie",
			tokenInCookie: true,
			tokenInHeader: false,
			invalidToken:  false,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Valid token in header",
			tokenInCookie: false,
			tokenInHeader: true,
			invalidToken:  false,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Valid token in cookie and invalid in header",
			tokenInCookie: true,
			tokenInHeader: false,
			invalidToken:  true,
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:          "Valid token in header and invalid in cookie",
			tokenInCookie: false,
			tokenInHeader: true,
			invalidToken:  true,
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:          "Invalid token in cookie and header",
			tokenInCookie: false,
			tokenInHeader: false,
			invalidToken:  true,
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:          "No token",
			tokenInCookie: false,
			tokenInHeader: false,
			invalidToken:  false,
			expectedCode:  http.StatusUnauthorized,
		},
	}

	john := []models.User{
		builders.NewUserBuilder().Build(),
	}

	config, _ := configs.LoadConfig()
	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, john[0].ID, config.AccessTokenPrivateKey)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	cookie := &http.Cookie{Name: "access_token", Value: accessToken}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRouter()
			router.GET("/", DeserializeUser(), func(context *gin.Context) {})
			defer sqlDB.Close()

			if tt.expectedCode == http.StatusOK {
				rows := testUtils.ConvertStructsToSQLMockRows(john)
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(john[0].ID, 1).
					WillReturnRows(rows)
			}
			req, _ := http.NewRequest("GET", "/", nil)

			if tt.tokenInCookie {
				if tt.invalidToken {
					req.AddCookie(&http.Cookie{Name: "access_token", Value: "invalidtoken"})
				} else {
					req.AddCookie(cookie)
				}
			}

			if tt.tokenInHeader {
				bearerToken := accessToken
				if tt.invalidToken {
					bearerToken = "invalidtoken"
				}
				req.Header.Add("Authorization", "Bearer "+bearerToken)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("SignUpInput.Validate() code = %v, expected code %v", w.Code, tt.expectedCode)
			}
		})
	}
}

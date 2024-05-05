package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/enzo-gbd/GBA/configs"
	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/enzo-gbd/GBA/internal/middlewares"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/models/builders"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/enzo-gbd/GBA/internal/utils/testUtils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var authController = NewAuthController()
var router *gin.Engine
var database *gorm.DB
var sqlDB *sql.DB
var mock sqlmock.Sqlmock

func setupRouter() {
	router = gin.Default()
	database, sqlDB, mock = db.InitMockDB()

	router.Use(middlewares.InjectDB(database))
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestSignUpInput(t *testing.T) {
	method, url := "POST", "/register"
	queryCreate := `INSERT INTO "users"`

	tests := []struct {
		name         string
		inputs       []models.SignUpInput
		expectedCode int
	}{
		{
			name: "valid input",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().BuildSignUpInput(),
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid FirstName",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WhereFirstName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid Name",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WhereName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid Birthday",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WhereBirthday(time.Time{}).BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid Gender",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WhereGender("none").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid Email",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WhereEmail("john.doe").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "short password",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WherePassword("Short1.").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "password without specials",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WherePassword("Password123").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "password without numbers",
			inputs: []models.SignUpInput{
				builders.NewUserBuilder().WherePassword("Password.").BuildSignUpInput(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "No body",
			inputs:       []models.SignUpInput{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRouter()
			router.POST(url, authController.SignUpUser)
			defer sqlDB.Close()

			if tt.expectedCode >= 200 && tt.expectedCode < 300 {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryCreate)).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			var model models.SignUpInput

			if len(tt.inputs) > 0 {
				model = tt.inputs[0]
			}
			w, err := utils.HttpTestRequest(router, method, url, model)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("code = %v, expected code %v", w.Code, tt.expectedCode)
			}
		})
	}
}

func TestSignInInput(t *testing.T) {
	method, url := "POST", "/login"
	queryFirst := `SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`

	tests := []struct {
		name          string
		input         models.SignInInput
		expectedError bool
		expectedCode  int
	}{
		{
			name:          "valid input",
			input:         builders.NewUserBuilder().BuildSignInInput(),
			expectedError: false,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "invalid credentials",
			input:         builders.NewUserBuilder().WherePassword("Password456.").BuildSignInInput(),
			expectedError: true,
			expectedCode:  http.StatusUnauthorized,
		},
		{
			name:          "invalid Email",
			input:         builders.NewUserBuilder().WhereEmail("john.doe").BuildSignInInput(),
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "short password",
			input:         builders.NewUserBuilder().WherePassword("Short1.").BuildSignInInput(),
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "password without specials",
			input:         builders.NewUserBuilder().WherePassword("Password123").BuildSignInInput(),
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "password without numbers",
			input:         builders.NewUserBuilder().WherePassword("Password.").BuildSignInInput(),
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "No body",
			input:         models.SignInInput{},
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
	}

	hashedPassword, err := utils.HashPassword("Password123.")
	if err != nil {
		t.Errorf("error = %v", err)
	}
	john := []models.User{
		builders.NewUserBuilder().WherePassword(hashedPassword).Build(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRouter()
			router.POST(url, authController.SignInUser)
			defer sqlDB.Close()

			if !tt.expectedError {
				rows := testUtils.ConvertStructsToSQLMockRows(john)
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(john[0].Email, 1).
					WillReturnRows(rows)
			}

			w, err := utils.HttpTestRequest(router, method, url, &tt.input)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("Scode = %v, expected code %v", w.Code, tt.expectedCode)
			}
			response := w.Result()
			cookiesNames := [3]string{"access_token", "refresh_token", "logged_in"}
			defer response.Body.Close()

			cookies := response.Cookies()

			for _, cookieName := range cookiesNames {
				var tokenCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == cookieName {
						tokenCookie = cookie
						break
					}
				}

				if (tokenCookie == nil) != tt.expectedError {
					t.Errorf("Expected %v cookie to be set", cookieName)
				} else if tokenCookie != nil == tt.expectedError {
					t.Errorf("Expected %v cookie to not be set", cookieName)
				} else if tokenCookie != nil {
					if tokenCookie.Value == "" {
						t.Errorf("Expected %v cookie to be not empty", cookieName)
					}
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
	method, url := "POST", "/logout"

	setupRouter()
	router.POST(url, authController.LogoutUser)
	defer sqlDB.Close()

	john := builders.NewUserBuilder().Build()

	config, _ := configs.LoadConfig()
	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, john.ID, config.AccessTokenPrivateKey)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	cookie := &http.Cookie{Name: "access_token", Value: accessToken}

	req, _ := http.NewRequest(method, url, nil)
	req.AddCookie(cookie)

	req.Header.Add("Authorization", "Bearer "+accessToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshAccessToken(t *testing.T) {
	method, url := "GET", "/refresh"
	queryFirst := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`
	tests := []struct {
		name                string
		hasRefreshToken     bool
		invalidRefreshToken bool
		expectedError       bool
		expectedCode        int
	}{
		{
			name:                "Valid refresh token",
			hasRefreshToken:     true,
			invalidRefreshToken: false,
			expectedError:       false,
			expectedCode:        http.StatusOK,
		},
		{
			name:                "Invalid refresh token",
			hasRefreshToken:     true,
			invalidRefreshToken: true,
			expectedError:       true,
			expectedCode:        http.StatusUnauthorized,
		},
		{
			name:                "No refresh token",
			hasRefreshToken:     false,
			invalidRefreshToken: false,
			expectedError:       true,
			expectedCode:        http.StatusUnauthorized,
		},
	}

	john := []models.User{
		builders.NewUserBuilder().Build(),
	}

	config, _ := configs.LoadConfig()
	refreshToken, err := utils.GenerateToken(config.RefreshTokenExpiresIn, john[0].ID, config.RefreshTokenPrivateKey)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	cookie := &http.Cookie{Name: "refresh_token", Value: refreshToken}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRouter()
			router.GET(url, authController.RefreshAccessToken)
			defer sqlDB.Close()

			if !tt.expectedError {
				rows := testUtils.ConvertStructsToSQLMockRows(john)
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(john[0].ID, 1).
					WillReturnRows(rows)
			}
			req, _ := http.NewRequest(method, url, nil)

			if tt.hasRefreshToken {
				if tt.invalidRefreshToken {
					req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalidtoken"})
				} else {
					req.AddCookie(cookie)
				}
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("code = %v, expected code %v", w.Code, tt.expectedCode)
			}
			response := w.Result()
			cookiesNames := [2]string{"access_token", "logged_in"}
			defer response.Body.Close()

			cookies := response.Cookies()

			for _, cookieName := range cookiesNames {
				var tokenCookie *http.Cookie
				for _, cookie := range cookies {
					if cookie.Name == cookieName {
						tokenCookie = cookie
						break
					}
				}

				if (tokenCookie == nil) != tt.expectedError {
					t.Errorf("Expected %v cookie to be set", cookieName)
				} else if tokenCookie != nil == tt.expectedError {
					t.Errorf("Expected %v cookie to not be set", cookieName)
				} else if tokenCookie != nil {
					if tokenCookie.Value == "" {
						t.Errorf("Expected %v cookie to be not empty", cookieName)
					}
				}
			}
		})
	}
}

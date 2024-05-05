package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/enzo-gbd/GBA/internal/middlewares"
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/models/builders"
	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/enzo-gbd/GBA/internal/utils/testUtils"
	"gorm.io/gorm"

	"github.com/enzo-gbd/GBA/internal/db"
	"github.com/gin-gonic/gin"
)

var userController = NewUserController()
var router *gin.Engine
var database *gorm.DB
var sqlDB *sql.DB
var mock sqlmock.Sqlmock

func SetupRouter() {
	router = gin.Default()
	database, sqlDB, mock = db.InitMockDB()

	router.Use(middlewares.InjectDB(database))
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestGetUsers(t *testing.T) {
	method, url := "GET", "/"
	query := "^SELECT (.+) FROM \"users\"$"

	tests := []struct {
		name         string
		items        []models.User
		expectedCode int
	}{
		{
			name: "One user",
			items: []models.User{
				builders.NewUserBuilder().Build(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Multiple users",
			items: []models.User{
				builders.NewUserBuilder().Build(),
				builders.NewUserBuilder().
					WhereFirstName("Julie").
					WhereGender("female").
					WhereEmail("julie.doe@mail.pe").
					Build(),
				builders.NewUserBuilder().
					WhereFirstName("Jeanne").
					WhereGender("female").
					WhereEmail("jeanne.doe@mail.pe").
					Build(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "No user",
			items:        []models.User{},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupRouter()
			router.GET(url, userController.GetUsers)
			defer sqlDB.Close()

			rows := testUtils.ConvertStructsToSQLMockRows(tt.items)
			mock.ExpectQuery(query).WillReturnRows(rows)

			w, err := utils.HttpTestRequest(router, method, url, nil)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var users []interface{}
				err = json.Unmarshal(w.Body.Bytes(), &users)
				if err != nil {
					t.Errorf("error unmarshalling response = %v", err)
				}

				if len(users) != len(tt.items) {
					t.Errorf("number of users = %v, expected %v", len(users), len(tt.items))
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	method, url := "GET", "/"
	query := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`

	john := []models.User{
		builders.NewUserBuilder().Build(),
	}

	tests := []struct {
		name          string
		idString      string
		expectedItems []models.User
		expectedCode  int
	}{
		{
			name:          "Valid id",
			idString:      john[0].ID.String(),
			expectedItems: john,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Invalid id but valid format",
			idString:      "cd1ef74e-4236-40fc-9542-614c03271cc7",
			expectedItems: nil,
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "Invalid id",
			idString:      "1234",
			expectedItems: nil,
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupRouter()
			router.GET(url+":id", userController.GetUserByID)
			defer sqlDB.Close()

			if tt.expectedItems != nil {
				rows := testUtils.ConvertStructsToSQLMockRows(tt.expectedItems)
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(tt.idString, 1).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(tt.idString, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			}
			w, err := utils.HttpTestRequest(router, method, url+tt.idString, nil)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var user *models.User
				err = json.Unmarshal(w.Body.Bytes(), &user)
				if err != nil {
					t.Errorf("error unmarshalling response = %v", err)
				}

				if user == nil {
					t.Errorf("user got is nil, expected not")
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	method, url := "PUT", "/"
	queryUpdate := `UPDATE "users" SET "first_name"=$1,"name"=$2,"birthday"=$3,"gender"=$4,"email"=$5,"password"=$6,"role"=$7,"address"=$8,"subscription_code"=$9,"is_active"=$10,"verification_code"=$11,"verified"=$12,"created_at"=$13,"updated_at"=$14,"deleted_at"=$15 WHERE "id" = $16`
	queryFirst := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`

	tests := []struct {
		name          string
		input         models.User
		idString      string
		expectedError bool
		expectedCode  int
	}{
		{
			name:          "valid input",
			input:         builders.NewUserBuilder().Build(),
			idString:      "",
			expectedError: false,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "invalid FirstName",
			input:         builders.NewUserBuilder().WhereFirstName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid Name",
			input:         builders.NewUserBuilder().WhereName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid Birthday",
			input:         builders.NewUserBuilder().WhereBirthday(time.Time{}).Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid Gender",
			input:         builders.NewUserBuilder().WhereGender("none").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "invalid Email",
			input:         builders.NewUserBuilder().WhereEmail("john.doe").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "short password",
			input:         builders.NewUserBuilder().WherePassword("Short1.").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "password without specials",
			input:         builders.NewUserBuilder().WherePassword("Password123").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "password without numbers",
			input:         builders.NewUserBuilder().WherePassword("Password.").Build(),
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "No body",
			input:         models.User{},
			idString:      "",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "Invalid id but valid format",
			input:         builders.NewUserBuilder().Build(),
			idString:      "cd1ef74e-4236-40fc-9542-614c03271cc7",
			expectedError: true,
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "Invalid id",
			input:         builders.NewUserBuilder().Build(),
			idString:      "1234",
			expectedError: true,
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupRouter()
			router.PUT(url+":id", userController.UpdateUser)
			defer sqlDB.Close()

			var idString string
			if tt.idString != "" {
				idString = tt.idString
			} else {
				idString = tt.input.ID.String()
			}
			if !tt.expectedError {
				rows := testUtils.ConvertStructsToSQLMockRows([]models.User{tt.input})
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(idString, 1).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
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
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(idString, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).
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
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectCommit()
			}
			w, err := utils.HttpTestRequest(router, method, url+idString, tt.input)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	method, url := "DELETE", "/"
	queryDelete := `DELETE FROM "users" WHERE (.+)`
	queryFirst := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`

	john := []models.User{
		builders.NewUserBuilder().Build(),
	}

	tests := []struct {
		name          string
		idString      string
		expectedItems []models.User
		expectedCode  int
	}{
		{
			name:          "Valid id",
			idString:      john[0].ID.String(),
			expectedItems: john,
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Invalid id but valid format",
			idString:      "cd1ef74e-4236-40fc-9542-614c03271cc7",
			expectedItems: nil,
			expectedCode:  http.StatusNotFound,
		},
		{
			name:          "Invalid id",
			idString:      "1234",
			expectedItems: nil,
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupRouter()
			router.DELETE(url+":id", userController.DeleteUser)
			defer sqlDB.Close()

			if tt.expectedItems != nil {
				rows := testUtils.ConvertStructsToSQLMockRows(tt.expectedItems)
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(tt.idString, 1).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec(queryDelete).
					WithArgs(tt.idString).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(queryFirst)).
					WithArgs(tt.idString, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectBegin()
				mock.ExpectExec(queryDelete).
					WithArgs(tt.idString).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectCommit()
			}

			w, err := utils.HttpTestRequest(router, method, url+tt.idString, nil)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}
		})
	}
}

func TestGetMe(t *testing.T) {
	method, url := "GET", "/me"

	john := []models.User{
		builders.NewUserBuilder().Build(),
	}

	tests := []struct {
		name         string
		currentUser  []models.User
		exists       bool
		expectedCode int
	}{
		{
			name:         "Valid user",
			currentUser:  john,
			exists:       true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "deleted user",
			currentUser:  john,
			exists:       false,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "no user",
			currentUser:  nil,
			exists:       false,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupRouter()
			router.GET(url, func(context *gin.Context) {
				if tt.exists {
					context.Set("currentUser", &tt.currentUser[0])
				}
			}, userController.GetMe)
			defer sqlDB.Close()

			w, err := utils.HttpTestRequest(router, method, url, nil)
			if err != nil {
				t.Errorf("error = %v", err)
			}
			if w.Code != tt.expectedCode {
				t.Errorf("HTTP status code = %v, expected %v", w.Code, tt.expectedCode)
			}

			if tt.expectedCode == http.StatusOK {
				var user *models.UserResponse
				err = json.Unmarshal(w.Body.Bytes(), &user)
				if err != nil {
					t.Errorf("error unmarshalling response = %v", err)
				}

				if user == nil {
					t.Errorf("user got is nil, expected not")
				}
			}
		})
	}
}

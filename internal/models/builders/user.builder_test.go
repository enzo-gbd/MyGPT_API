package builders

import (
	"database/sql"
	"testing"
	"time"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/google/uuid"
)

func TestUserBuilder(t *testing.T) {
	now := time.Now()
	john := models.User{
		ID:        uuid.New(),
		FirstName: "John",
		Name:      "Doe",
		Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		Gender:    "male",
		Email:     "john.doe@mail.pe",
		Password:  "Password123.",
		Role:      "user",
		Address:   sql.NullString{String: "123 Main St", Valid: true},
	}

	tests := []struct {
		name         string
		method       func() models.User
		param        interface{}
		expectedUser models.User
	}{
		{
			name:   "WithFirstName",
			method: NewUserBuilder().WithBase(john).WhereFirstName("Julie").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "Julie",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithName",
			method: NewUserBuilder().WithBase(john).WhereName("Does").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Does",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithBirthday",
			method: NewUserBuilder().WithBase(john).WhereBirthday(time.Date(2002, time.November, 2, 0, 0, 0, 0, time.UTC)).Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(2002, time.November, 2, 0, 0, 0, 0, time.UTC),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithGender",
			method: NewUserBuilder().WithBase(john).WhereGender("female").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "female",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithEmail",
			method: NewUserBuilder().WithBase(john).WhereEmail("julie.does@mail.pe").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "julie.does@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithPassword",
			method: NewUserBuilder().WithBase(john).WherePassword("Password456!").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password456!",
				Role:      "user",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithRole",
			method: NewUserBuilder().WithBase(john).WhereRole("admin").Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "admin",
				Address:   sql.NullString{String: "123 Main St", Valid: true},
			},
		},
		{
			name:   "WithAddress",
			method: NewUserBuilder().WithBase(john).WhereAddress(sql.NullString{String: "456 Main St", Valid: true}).Build,
			expectedUser: models.User{
				ID:        uuid.New(),
				FirstName: "John",
				Name:      "Doe",
				Birthday:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
				Gender:    "male",
				Email:     "john.doe@mail.pe",
				Password:  "Password123.",
				Role:      "user",
				Address:   sql.NullString{String: "456 Main St", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectedUser.IsEqualTo(tt.method()) {
				t.Errorf("%v expected = %v\ngot %v", tt.name, tt.expectedUser, tt.method())
			}
		})
	}
}

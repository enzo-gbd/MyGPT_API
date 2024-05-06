package models_test

import (
	"testing"
	"time"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/enzo-gbd/GBA/internal/models/builders"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name          string
		input         models.User
		expectedError bool
	}{
		{
			name:          "valid input",
			input:         builders.NewUserBuilder().Build(),
			expectedError: false,
		},
		{
			name:          "invalid FirstName",
			input:         builders.NewUserBuilder().WhereFirstName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").Build(),
			expectedError: true,
		},
		{
			name:          "invalid Name",
			input:         builders.NewUserBuilder().WhereName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").Build(),
			expectedError: true,
		},
		{
			name:          "invalid Birthday",
			input:         builders.NewUserBuilder().WhereBirthday(time.Time{}).Build(),
			expectedError: true,
		},
		{
			name:          "invalid Gender",
			input:         builders.NewUserBuilder().WhereGender("none").Build(),
			expectedError: true,
		},
		{
			name:          "invalid Email",
			input:         builders.NewUserBuilder().WhereEmail("john.doe").Build(),
			expectedError: true,
		},
		{
			name:          "short password",
			input:         builders.NewUserBuilder().WherePassword("Short1.").Build(),
			expectedError: true,
		},
		{
			name:          "password without specials",
			input:         builders.NewUserBuilder().WherePassword("Password123").Build(),
			expectedError: true,
		},
		{
			name:          "password without numbers",
			input:         builders.NewUserBuilder().WherePassword("Password.").Build(),
			expectedError: true,
		},
		{
			name:          "invalid Role",
			input:         builders.NewUserBuilder().WhereRole("none").Build(),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("SignUpInput.Validate() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	user := &models.User{}
	err := user.BeforeCreate(nil)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, user.ID)
}

func TestSignUpInputValidation(t *testing.T) {
	tests := []struct {
		name          string
		input         models.SignUpInput
		expectedError bool
	}{
		{
			name:          "valid input",
			input:         builders.NewUserBuilder().BuildSignUpInput(),
			expectedError: false,
		},
		{
			name:          "invalid FirstName",
			input:         builders.NewUserBuilder().WhereFirstName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "invalid Name",
			input:         builders.NewUserBuilder().WhereName("ncjebvcizbckzbclozbcozbcmabckecaveaveaec").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "invalid Birthday",
			input:         builders.NewUserBuilder().WhereBirthday(time.Time{}).BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "invalid Gender",
			input:         builders.NewUserBuilder().WhereGender("none").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "invalid Email",
			input:         builders.NewUserBuilder().WhereEmail("john.doe").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "short password",
			input:         builders.NewUserBuilder().WherePassword("Short1.").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "password without specials",
			input:         builders.NewUserBuilder().WherePassword("Password123").BuildSignUpInput(),
			expectedError: true,
		},
		{
			name:          "password without numbers",
			input:         builders.NewUserBuilder().WherePassword("Password.").BuildSignUpInput(),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("SignUpInput.Validate() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestSignInInputValidation(t *testing.T) {
	tests := []struct {
		name          string
		input         models.SignInInput
		expectedError bool
	}{
		{
			name:          "valid input",
			input:         builders.NewUserBuilder().BuildSignInInput(),
			expectedError: false,
		},
		{
			name:          "invalid Email",
			input:         builders.NewUserBuilder().WhereEmail("john.doe").BuildSignInInput(),
			expectedError: true,
		},
		{
			name:          "short password",
			input:         builders.NewUserBuilder().WherePassword("Short1.").BuildSignInInput(),
			expectedError: true,
		},
		{
			name:          "password without specials",
			input:         builders.NewUserBuilder().WherePassword("Password123").BuildSignInInput(),
			expectedError: true,
		},
		{
			name:          "password without numbers",
			input:         builders.NewUserBuilder().WherePassword("Password.").BuildSignInInput(),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("SignUpInput.Validate() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

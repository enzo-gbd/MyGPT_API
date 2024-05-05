package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/enzo-gbd/GBA/internal/utils"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

// User represents a user profile in the system.
// @Description User holds all details about a user.
type User struct {
	ID               uuid.UUID      `gorm:"type:char(36);primary_key" `             // Unique identifier for the user
	FirstName        string         `gorm:"type:varchar(255);not null"`             // First name of the user
	Name             string         `gorm:"type:varchar(255);not null"`             // Last name of the user
	Birthday         time.Time      `gorm:"not null"`                               // Birthday of the user
	Gender           string         `gorm:"type:varchar(255);not null"`             // Gender of the user
	Email            string         `gorm:"type:varchar(255);uniqueIndex;not null"` // Email address of the user, must be unique
	Password         string         `gorm:"type:varchar(255);not null"`             // Password for the user account
	Role             string         `gorm:"type:varchar(255);default:user"`         // Role of the user in the system
	Address          sql.NullString `gorm:"type:varchar(255)"`                      // Optional address of the user
	SubscriptionCode sql.NullString `gorm:"type:varchar(255)"`                      // Optional subscription code
	IsActive         bool           `gorm:"default:1"`                              // Flag indicating if the user account is active
	VerificationCode sql.NullString // Optional verification code
	Verified         bool           `gorm:"not null;default:0"` // Flag indicating if the user has verified their email
	CreatedAt        time.Time      `gorm:"not null"`           // Timestamp when the user was created
	UpdatedAt        time.Time      `gorm:"not null"`           // Timestamp when the user was last updated
	DeletedAt        time.Time      // Optional timestamp when the user was deleted
}

// Validate performs validation on User fields using ozzo-validation package.
// It ensures all necessary fields meet their expected format.
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.Birthday, validation.Required),
		validation.Field(&u.Gender, validation.Required, validation.In("male", "female", "other")),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100), is.PrintableASCII, validation.By(utils.PasswordRequirements)),
		validation.Field(&u.Role, validation.Required, validation.In("user", "admin")),
	)
}

// BeforeCreate is a GORM hook that is called before a new user record is created.
// It assigns a new UUID to the user's ID.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// String provides a string representation of the user which includes
// personal details and contact information.
func (u User) String() string {
	return fmt.Sprintf("%v %v (%v)(%v)\nemail: %v\nadress: %v", u.FirstName, u.Name, u.Birthday, u.Gender, u.Email, u.Address.String)
}

// IsEqualTo compares this User instance with another to determine if they are identical.
func (u User) IsEqualTo(other User) bool {
	return u.FirstName == other.FirstName &&
		u.Name == other.Name &&
		u.Birthday == other.Birthday &&
		u.Gender == other.Gender &&
		u.Email == other.Email &&
		u.Password == other.Password &&
		u.Role == other.Role &&
		u.Address == other.Address
}

// SignUpInput represents the required fields for a user register.
// @Description Fields required to register a new user.
type SignUpInput struct {
	FirstName string    `json:"first_name" binding:"required"` // First name of the user
	Name      string    `json:"name" binding:"required"`       // Last name of the user
	Birthday  time.Time `json:"birthday" binding:"required"`   // Birthday of the user
	Gender    string    `json:"gender" binding:"required"`     // Gender of the user
	Email     string    `json:"email" binding:"required"`      // Email address of the user
	Password  string    `json:"password" binding:"required"`   // Password for the user account
}

// Validate performs validation on SignUpInput fields to ensure they meet
// the requirements for user registration.
func (s SignUpInput) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FirstName, validation.Required, validation.Length(1, 20)),
		validation.Field(&s.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&s.Birthday, validation.Required),
		validation.Field(&s.Gender, validation.Required, validation.In("male", "female", "other")),
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 100), is.PrintableASCII, validation.By(utils.PasswordRequirements)),
	)
}

// SignInInput represents the required fields for a user login.
// @Description Fields required to login a user.
type SignInInput struct {
	Email    string `json:"email" binding:"required"`    // Email address of the user
	Password string `json:"password" binding:"required"` // Password for the user account
}

// Validate performs validation on SignInInput fields to ensure they meet
// the requirements for user authentication.
func (s SignInInput) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 100), is.PrintableASCII, validation.By(utils.PasswordRequirements)),
	)
}

// UserResponse represents the public user profile that is returned by the API.
// @Description UserResponse holds the data that is exposed to the client after API requests.
type UserResponse struct {
	ID               uuid.UUID `json:"id,omitempty"`         // Unique identifier for the user, omitted if empty
	FirstName        string    `json:"first_name,omitempty"` // First name of the user, omitted if empty
	Name             string    `json:"name,omitempty"`       // Last name of the user, omitted if empty
	Birthday         time.Time `json:"birthday,omitempty"`   // Birthday of the user, omitted if empty
	Gender           string    `json:"gender,omitempty"`     // Gender of the user, omitted if empty
	Email            string    `json:"email,omitempty"`      // Email address of the user, omitted if empty
	Role             string    `json:"role,omitempty"`       // Role of the user within the system, omitted if empty
	Address          string    `json:"address"`              // Address of the user (present even if empty as a blank string)
	SubscriptionCode string    `json:"subscription_code"`    // Subscription code related to the user's account (present even if empty as a blank string)
	CreatedAt        time.Time `json:"created_at"`           // Timestamp when the user was created
	UpdatedAt        time.Time `json:"updated_at"`           // Timestamp when the user profile was last updated
}

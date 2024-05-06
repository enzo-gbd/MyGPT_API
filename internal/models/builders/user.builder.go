// Package builders provides the functionality to construct and configure models
// with various properties set, using the builder pattern.
package builders

import (
	"database/sql"
	"time"

	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/google/uuid"
)

// UserBuilder is a builder used to construct a user model with various properties.
type UserBuilder struct {
	u models.User
}

// NewUserBuilder initializes and returns a new UserBuilder with default values for a User.
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		u: models.User{
			ID:        uuid.New(),
			FirstName: "John",
			Name:      "Doe",
			Birthday:  time.Now(),
			Gender:    "male",
			Email:     "john.doe@mail.pe",
			Password:  "Password123.",
			Role:      "user",
			Address:   sql.NullString{String: "123 Main St", Valid: true},
		},
	}
}

// Build constructs and returns the User model configured by the builder.
func (ub *UserBuilder) Build() models.User {
	return ub.u
}

// BuildSignUpInput constructs and returns a SignUpInput model based on the User model configured by the builder.
func (ub *UserBuilder) BuildSignUpInput() models.SignUpInput {
	return models.SignUpInput{
		FirstName: ub.u.FirstName,
		Name:      ub.u.Name,
		Birthday:  ub.u.Birthday,
		Gender:    ub.u.Gender,
		Email:     ub.u.Email,
		Password:  ub.u.Password,
	}
}

// BuildSignInInput constructs and returns a SignInInput model using the email and password from the User model.
func (ub *UserBuilder) BuildSignInInput() models.SignInInput {
	return models.SignInInput{
		Email:    ub.u.Email,
		Password: ub.u.Password,
	}
}

// WithBase sets the base user model and returns the UserBuilder.
func (ub *UserBuilder) WithBase(user models.User) *UserBuilder {
	ub.u = user
	return ub
}

// WhereFirstName sets the FirstName of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereFirstName(firstName string) *UserBuilder {
	ub.u.FirstName = firstName
	return ub
}

// WhereName sets the Name of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereName(name string) *UserBuilder {
	ub.u.Name = name
	return ub
}

// WhereBirthday sets the Birthday of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereBirthday(birthday time.Time) *UserBuilder {
	ub.u.Birthday = birthday
	return ub
}

// WhereGender sets the Gender of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereGender(gender string) *UserBuilder {
	ub.u.Gender = gender
	return ub
}

// WhereEmail sets the Email of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereEmail(email string) *UserBuilder {
	ub.u.Email = email
	return ub
}

// WherePassword sets the Password of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WherePassword(password string) *UserBuilder {
	ub.u.Password = password
	return ub
}

// WhereRole sets the Role of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereRole(role string) *UserBuilder {
	ub.u.Role = role
	return ub
}

// WhereAddress sets the Address of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereAddress(address sql.NullString) *UserBuilder {
	ub.u.Address = address
	return ub
}

// WhereSubscriptionCode sets the SubscriptionCode of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereSubscriptionCode(subscriptionCode sql.NullString) *UserBuilder {
	ub.u.SubscriptionCode = subscriptionCode
	return ub
}

// WhereIsActive sets the IsActive flag of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereIsActive(isActive bool) *UserBuilder {
	ub.u.IsActive = isActive
	return ub
}

// WhereVerificationCode sets the VerificationCode of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereVerificationCode(verificationCode sql.NullString) *UserBuilder {
	ub.u.VerificationCode = verificationCode
	return ub
}

// WhereCreatedAt sets the CreatedAt timestamp of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereCreatedAt(createdAt time.Time) *UserBuilder {
	ub.u.CreatedAt = createdAt
	return ub
}

// WhereUpdatedAt sets the UpdatedAt timestamp of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereUpdatedAt(updatedAt time.Time) *UserBuilder {
	ub.u.UpdatedAt = updatedAt
	return ub
}

// WhereDeletedAt sets the DeletedAt timestamp of the user being built and returns the UserBuilder.
func (ub *UserBuilder) WhereDeletedAt(deletedAt time.Time) *UserBuilder {
	ub.u.DeletedAt = deletedAt
	return ub
}

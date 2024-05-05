// Package utils provides utility functions that support various operations across the application.
package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain text password and returns a bcrypt hashed version of it.
// If hashing fails, an error is returned. The bcrypt.DefaultCost is used for hashing.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a bcrypt hashed password with a candidate password.
// If the passwords do not match, an error is returned, indicating the verification failure.
func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

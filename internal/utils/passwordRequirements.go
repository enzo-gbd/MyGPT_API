// Package utils provides utility functions that support various operations across the application.
package utils

import (
	"errors"
	"strings"
)

// PasswordRequirements checks if the given password meets specific security requirements.
// The function takes a value of any type, but expects a string for validation.
// It validates the password based on the following criteria:
// - Must contain at least one digit.
// - Must contain at least one special character.
// - Must contain at least one uppercase letter.
// - Must contain at least one lowercase letter.
// If the password does not meet these criteria, it returns an error listing all failed requirements.
// If the password meets all the criteria, it returns nil indicating no error.
func PasswordRequirements(value interface{}) error {
	s, _ := value.(string) // attempts to cast value to a string
	var errMessages []string

	// Check for at least one digit
	if !strings.ContainsAny(s, "0123456789") {
		errMessages = append(errMessages, "must contain at least one digit")
	}

	// Check for at least one special character
	if !strings.ContainsAny(s, "*!.@$%^&(){}[]:;<>,.?/~_+-=") {
		errMessages = append(errMessages, "must contain at least one special character")
	}

	// Check for at least one uppercase letter
	if !strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		errMessages = append(errMessages, "must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyz") {
		errMessages = append(errMessages, "must contain at least one lowercase letter")
	}

	// If there are error messages, join them and return as an error
	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, ", "))
	}

	// If all checks are passed, return nil
	return nil
}

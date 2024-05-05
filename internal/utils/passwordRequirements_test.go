package utils

import (
	"strings"
	"testing"
)

func TestPasswordRequirements(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
		errMsg   string
	}{
		{"Valid1@", false, ""},
		{"NoDigit", true, "must contain at least one digit"},
		{"Nodigitorspecial@", true, "must contain at least one digit"},
		{"NODIGITUPPERCASE@", true, "must contain at least one digit, must contain at least one lowercase letter"},
		{"12345", true, "must contain at least one special character, must contain at least one uppercase letter, must contain at least one lowercase letter"},
		{"alllowercase12345", true, "must contain at least one special character, must contain at least one uppercase letter"},
	}

	for _, tc := range tests {
		t.Run(tc.password, func(t *testing.T) {
			err := PasswordRequirements(tc.password)
			if tc.wantErr && err == nil {
				t.Errorf("Expected an error but got none")
			} else if !tc.wantErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("Expected error message to contain '%v', got '%v'", tc.errMsg, err.Error())
			}
		})
	}
}

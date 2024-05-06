package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "secretPassword"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	require.NotEqual(t, password, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	password := "secretPassword"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = VerifyPassword(hashedPassword, password)
	require.NoError(t, err)

	wrongPassword := "wrongPassword"
	err = VerifyPassword(hashedPassword, wrongPassword)
	require.Error(t, err)
}

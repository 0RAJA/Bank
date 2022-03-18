package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(10)

	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotZero(t, hashPassword)
	require.NoError(t, CheckPassword(password, hashPassword))
	wrongPassword := RandomString(10)
	require.Error(t, CheckPassword(wrongPassword, hashPassword))
}

package db

import (
	"context"
	"github.com/0RAJA/Bank/db/util"
	"github.com/0RAJA/Bank/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueries_CreateUser(t *testing.T) {
	hashPassword, err := utils.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		Email:          util.RandomString(10),
		FullName:       util.RandomOwner(),
	}
	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       arg.Username,
		HashedPassword: arg.HashedPassword,
		Email:          arg.Email,
		FullName:       arg.FullName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.Zero(t, user.PasswordChangedAt)
}
func testCreateUser(t *testing.T) User {
	hashPassword, err := utils.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		Email:          util.RandomString(10),
		FullName:       util.RandomOwner(),
	}
	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       arg.Username,
		HashedPassword: arg.HashedPassword,
		Email:          arg.Email,
		FullName:       arg.FullName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.Zero(t, user.PasswordChangedAt)
	return user
}

func TestQueries_GetUser(t *testing.T) {
	user := testCreateUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotZero(t, user2)

	require.Equal(t, user, user2)
}

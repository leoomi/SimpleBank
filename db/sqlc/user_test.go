package db

import (
	"context"
	"testing"

	"github.com/leoomi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	params := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.Password, user.Password)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), newUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, newUser.Username, user.Username)
	require.Equal(t, newUser.Password, user.Password)
	require.Equal(t, newUser.FullName, user.FullName)
	require.Equal(t, newUser.Email, user.Email)
	require.Equal(t, newUser.CreatedAt, user.CreatedAt)
	require.Equal(t, newUser.PasswordChangedAt, user.PasswordChangedAt)
}

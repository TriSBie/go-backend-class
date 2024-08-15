package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"simple_bank.sqlc.dev/app/util"
)

func createRandomUser(t *testing.T) User {
	ctx := context.Background()

	hashedPassword, err := util.HashPassword(util.RandomString(6))

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	// using context with testQueries (*Queries) initialized from main test
	user, err := testQueries.CreateUser(ctx, arg)

	// Ensure the execution run without nil or error
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Email, arg.Email)

	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	// IsZero reports whether t represents the zero time instant
	require.True(t, user.PasswordChangedAt.IsZero())
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)

	// ensure two time of assert is within a delta duration
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

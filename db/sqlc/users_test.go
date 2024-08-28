package db

import (
	"context"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := RegisterUserParams{
		ID:       util.GenUUID(),
		Email:    util.RandomEmail(),
		Password: util.RandomPassword(),
	}

	user, err := testStore.RegisterUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestRegisterUser(t *testing.T) {
	createRandomUser(t)
}

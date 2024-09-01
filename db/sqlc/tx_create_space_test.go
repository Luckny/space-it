package db

import (
	"context"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSpaceTx(t *testing.T) {
	user := createRandomUser(t)

	n := 5

	errs := make(chan error)
	results := make(chan CreateSpaceTxResult)

	// run n concurrent create Space transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.CreateSpaceTx(context.Background(), CreateSpaceTxParams{
				Name:  util.RandomSpaceName(),
				Owner: user.ID,
			})

			errs <- err
			results <- result
		}()
	}

	var space Space
	var permission Permission

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check space
		space = result.Space
		require.NotEmpty(t, space)
		require.Equal(t, user.ID, space.Owner)
		require.NotZero(t, space.ID)
		require.NotZero(t, space.Name)
		require.NotZero(t, space.CreatedAt)

		// check permission
		permission = result.Permission
		require.NotEmpty(t, permission)
		require.Equal(t, permission.UserID, user.ID)
		require.NotZero(t, permission.SpaceID)
		require.True(t, permission.ReadPermission)
		require.True(t, permission.WritePermission)
		require.True(t, permission.DeletePermission)
		require.NotZero(t, permission.CreatedAt)
		require.NotZero(t, permission.UpdatedAt)

	}

	_, err := testStore.GetSpaceByID(context.Background(), space.ID)
	require.NoError(t, err)

	getPermArg := GetPermissionsByUserAndSpaceIDParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}
	_, err = testStore.GetPermissionsByUserAndSpaceID(context.Background(), getPermArg)
	require.NoError(t, err)
}

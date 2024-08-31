package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateAdminPermission(t *testing.T) {
	user := createRandomUser(t)
	space := createRandomSpace(t, user)
	createTestAdminPermission(t, user, space)
	user = createRandomUser(t)
	space = createRandomSpace(t, user)
	createTestReadPermission(t, user, space)
	user = createRandomUser(t)
	space = createRandomSpace(t, user)
	createTestWritePermission(t, user, space)
	user = createRandomUser(t)
	space = createRandomSpace(t, user)
	createTestDeletePermission(t, user, space)
	user = createRandomUser(t)
	space = createRandomSpace(t, user)
	createTestPermission(t, user, space)

}

func TestGetPermissionByUserAndSpaceId(t *testing.T) {
	user := createRandomUser(t)
	space := createRandomSpace(t, user)
	perm1 := createTestDeletePermission(t, user, space)

	arg := GetPermissionsByUserAndSpaceIDParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}

	perm2, err := testStore.GetPermissionsByUserAndSpaceID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, perm2)

	require.Equal(t, perm1.UserID, perm2.UserID)
	require.Equal(t, perm1.WritePermission, perm2.WritePermission)
	require.Equal(t, perm1.ReadPermission, perm2.ReadPermission)
	require.Equal(t, perm1.DeletePermission, perm2.DeletePermission)

	require.WithinDuration(t, perm1.CreatedAt.Time, perm2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, perm1.UpdatedAt.Time, perm2.UpdatedAt.Time, time.Second)
}

// ------ ------------ Helpers

func createTestAdminPermission(t *testing.T, user User, space Space) Permission {
	arg := CreateAllPermissionParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}

	permission, err := testStore.CreateAllPermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, permission.UserID)
	require.Equal(t, space.ID, permission.SpaceID)
	require.True(t, permission.ReadPermission)
	require.True(t, permission.WritePermission)
	require.True(t, permission.DeletePermission)

	require.NotZero(t, permission.CreatedAt)
	require.NotZero(t, permission.UpdatedAt)

	return permission
}

func createTestReadPermission(t *testing.T, user User, space Space) Permission {
	arg := CreateReadPermissionParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}

	permission, err := testStore.CreateReadPermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, permission.UserID)
	require.Equal(t, space.ID, permission.SpaceID)
	require.True(t, permission.ReadPermission)
	require.False(t, permission.WritePermission)
	require.False(t, permission.DeletePermission)

	require.NotZero(t, permission.CreatedAt)
	require.NotZero(t, permission.UpdatedAt)

	return permission
}

func createTestWritePermission(t *testing.T, user User, space Space) Permission {
	arg := CreateWritePermissionParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}

	permission, err := testStore.CreateWritePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, permission.UserID)
	require.Equal(t, space.ID, permission.SpaceID)
	require.False(t, permission.ReadPermission)
	require.True(t, permission.WritePermission)
	require.False(t, permission.DeletePermission)

	require.NotZero(t, permission.CreatedAt)
	require.NotZero(t, permission.UpdatedAt)

	return permission
}

func createTestDeletePermission(t *testing.T, user User, space Space) Permission {
	arg := CreateDeletePermissionParams{
		UserID:  user.ID,
		SpaceID: space.ID,
	}

	permission, err := testStore.CreateDeletePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, permission.UserID)
	require.Equal(t, space.ID, permission.SpaceID)
	require.False(t, permission.ReadPermission)
	require.False(t, permission.WritePermission)
	require.True(t, permission.DeletePermission)

	require.NotZero(t, permission.CreatedAt)
	require.NotZero(t, permission.UpdatedAt)

	return permission
}

func createTestPermission(t *testing.T, user User, space Space) Permission {
	arg := CreatePermissionParams{
		UserID:           user.ID,
		SpaceID:          space.ID,
		ReadPermission:   true,
		WritePermission:  true,
		DeletePermission: false,
	}

	permission, err := testStore.CreatePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, permission.UserID)
	require.Equal(t, space.ID, permission.SpaceID)

	require.Equal(t, arg.ReadPermission, permission.ReadPermission)
	require.Equal(t, arg.WritePermission, permission.WritePermission)
	require.Equal(t, arg.DeletePermission, permission.DeletePermission)

	require.NotZero(t, permission.CreatedAt)
	require.NotZero(t, permission.UpdatedAt)

	return permission
}

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Luckny/space-it/util"
	"github.com/stretchr/testify/require"
)

func createRandomSpace(t *testing.T, user User) Space {

	arg := CreateSpaceParams{
		Name:  util.RandomSpaceName(),
		Owner: user.ID,
	}

	space, err := testStore.CreateSpace(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, space)

	require.Equal(t, user.ID, space.Owner)
	require.Equal(t, arg.Name, space.Name)

	require.NotZero(t, space.ID)
	require.NotZero(t, space.CreatedAt)

	return space
}

func TestCreateSpace(t *testing.T) {
	user := createRandomUser(t)
	createRandomSpace(t, user)
}

func TestGetSpaceByID(t *testing.T) {
	user := createRandomUser(t)
	space1 := createRandomSpace(t, user)

	space2, err := testStore.GetSpaceByID(context.Background(), space1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, space2)

	require.Equal(t, space1.Name, space2.Name)
	require.Equal(t, space1.Owner, space2.Owner)
	require.Equal(t, space1.ID, space2.ID)

	require.WithinDuration(t, space1.CreatedAt.Time, space2.CreatedAt.Time, time.Second)
}

func TestGetSpaceByName(t *testing.T) {
	user := createRandomUser(t)
	space1 := createRandomSpace(t, user)

	space2, err := testStore.GetSpaceByName(context.Background(), space1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, space2)

	require.Equal(t, space1.Name, space2.Name)
	require.Equal(t, space1.Owner, space2.Owner)
	require.Equal(t, space1.ID, space2.ID)

	require.WithinDuration(t, space1.CreatedAt.Time, space2.CreatedAt.Time, time.Second)

}

func TestListSpaces(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomSpace(t, user)
	}

	arg := ListSpacesParams{
		Limit:  5,
		Offset: 5,
	}

	spaces, err := testStore.ListSpaces(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, spaces, 5)

	for _, space := range spaces {
		require.NotEmpty(t, space)
	}
}

func TestUpdateSpace(t *testing.T) {
	user := createRandomUser(t)
	space1 := createRandomSpace(t, user)

	arg := UpdateSpaceParams{
		ID:   space1.ID,
		Name: util.RandomSpaceName(),
	}

	space2, err := testStore.UpdateSpace(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, space2.ID, space1.ID)
	require.Equal(t, space2.Owner, space1.Owner)
	require.NotEqual(t, space2.Name, space1.Name)
	require.Equal(t, space2.Name, arg.Name)

	require.WithinDuration(t, space1.CreatedAt.Time, space2.CreatedAt.Time, time.Second)
}

func TestDeleteSpace(t *testing.T) {
	user := createRandomUser(t)
	space1 := createRandomSpace(t, user)
	err := testStore.DeleteSpace(context.Background(), space1.ID)
	require.NoError(t, err)

	space2, err := testStore.GetSpaceByID(context.Background(), space1.ID)
	require.Error(t, err)
	require.ErrorContains(t, sql.ErrNoRows, err.Error())
	require.Empty(t, space2)
}

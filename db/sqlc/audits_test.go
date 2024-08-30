package db

import (
	"context"
	"net/http"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createTestUnAuthenticatedRequestLog(t *testing.T, userID uuid.UUID) uuid.UUID {

	arg := CreateUnauthenticatedRequestLogParams{
		ID:     util.GenUUID(),
		Method: http.MethodGet,
		Path:   "/somepath",
	}

	reqLogId, err := testStore.CreateUnauthenticatedRequestLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, reqLogId)

	require.Equal(t, arg.ID, reqLogId)
	return reqLogId
}

func TestCreateRequestLog(t *testing.T) {
	user := createRandomUser(t)
	createTestUnAuthenticatedRequestLog(t, user.ID)
}

func TestCreateResponseLog(t *testing.T) {
	user := createRandomUser(t)
	reqLogId := createTestUnAuthenticatedRequestLog(t, user.ID)

	arg := CreateResponseLogParams{
		ID:     reqLogId,
		Status: http.StatusOK,
	}

	err := testStore.CreateResponseLog(context.Background(), arg)
	require.NoError(t, err)

}

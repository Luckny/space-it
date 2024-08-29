package db

import (
	"context"
	"net/http"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createTestRequestLog(t *testing.T, userID uuid.UUID) uuid.UUID {

	arg := CreateRequestLogParams{
		ID:     util.GenUUID(),
		Method: http.MethodGet,
		Path:   "/somepath",
		UserID: userID,
	}

	reqLogId, err := testStore.CreateRequestLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, reqLogId)

	require.Equal(t, arg.ID, reqLogId)
	return reqLogId
}

func TestCreateRequestLog(t *testing.T) {
	user := createRandomUser(t)
	createTestRequestLog(t, user.ID)
}

func TestCreateResponseLog(t *testing.T) {
	user := createRandomUser(t)
	reqLogId := createTestRequestLog(t, user.ID)

	arg := CreateResponseLogParams{
		ID:     reqLogId,
		Status: http.StatusOK,
	}

	err := testStore.CreateResponseLog(context.Background(), arg)
	require.NoError(t, err)

}

package db

import (
	"context"
	"net/http"
	"testing"

	"github.com/Luckny/space-it/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createTestUnAuthenticatedRequestLog(t *testing.T, userID uuid.UUID) RequestLog {

	arg := CreateUnauthenticatedRequestLogParams{
		ID:     util.GenUUID(),
		Method: http.MethodGet,
		Path:   "/somepath",
	}

	reqLog, err := testStore.CreateUnauthenticatedRequestLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, reqLog)

	require.Equal(t, arg.ID, reqLog.ID)
	require.Equal(t, arg.Method, reqLog.Method)
	require.Equal(t, arg.Path, reqLog.Path)

	require.NotZero(t, reqLog.CreatedAt)

	return reqLog
}

func TestCreateRequestLog(t *testing.T) {
	user := createRandomUser(t)
	createTestUnAuthenticatedRequestLog(t, user.ID)
}

func TestCreateResponseLog(t *testing.T) {
	user := createRandomUser(t)
	reqLog := createTestUnAuthenticatedRequestLog(t, user.ID)

	arg := CreateResponseLogParams{
		ID:     reqLog.ID,
		Status: http.StatusOK,
	}

	resLog, err := testStore.CreateResponseLog(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, resLog)

	require.Equal(t, arg.ID, resLog.ID, reqLog.ID)
	require.Equal(t, arg.Status, resLog.Status)

	require.NotZero(t, resLog.CreatedAt)
}

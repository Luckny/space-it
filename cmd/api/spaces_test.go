package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Luckny/space-it/db/mock"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateSpaceAPI(t *testing.T) {
	user, _ := mockdb.RandomUser(t)
	spaceTxResult := mockdb.RandomSpaceTxResult(t, user.ID)
	space := spaceTxResult.Space

	testCases := []struct {
		name          string
		body          createSpaceRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{

			name: "should create space",
			body: createSpaceRequest{
				Name: space.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateSpaceTxParams{
					Name:  space.Name,
					Owner: user.ID,
				}
				store.EXPECT().
					CreateSpaceTx(gomock.Any(), mockdb.EqCreateSpaceTxParam(arg)).
					Times(1).
					Return(spaceTxResult, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusCreated)
				requireBodyMatchSpace(t, recorder.Body, space)
			},
		},

		{
			name: "bad request",
			body: createSpaceRequest{
				Name: "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSpace(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},

		{
			name: "internal error",
			body: createSpaceRequest{
				Name: space.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateSpaceTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateSpaceTxResult{}, db.ErrForeignKeyConstraint)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		// Run test case
		t.Run(tc.name, func(t *testing.T) {

			// init gomock
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// api server with mock store
			server := NewServer(store, config.Config{})
			router := gin.Default()
			// middleware to add user in context if test is authenticated
			router.Use(func(ctx *gin.Context) {
				ctx.Set("user", &user)
				ctx.Next()
			})

			router.POST("/spaces", server.createSpace)

			// request params
			jsonBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// create request
			request, err := http.NewRequest(http.MethodPost, "/spaces", bytes.NewReader(jsonBody))
			require.NoError(t, err)

			// test recorder
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(recorder)
		})
	}
}

// requireBodyMatchUser checks that the user in the body matches the recieved user
func requireBodyMatchSpace(t *testing.T, body *bytes.Buffer, space db.Space) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotSpace db.Space
	err = json.Unmarshal(data, &gotSpace)
	gotSpace.ID = space.ID
	require.NoError(t, err)
	require.Equal(t, space.Name, gotSpace.Name)
	require.Equal(t, space.Owner, gotSpace.Owner)
}

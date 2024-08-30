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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserAPI(t *testing.T) {
	user, unHashedPassword := mockdb.RandomUser(t)

	// TODO: authenticated tests should call GetUserByEmail and CreateAuthenticatedRequestLog

	testCases := []struct {
		name          string
		body          registerUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "should register user",
			body: registerUserRequest{
				Email:    user.Email,
				Password: unHashedPassword,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.RegisterUserParams{
					ID:       user.ID,
					Email:    user.Email,
					Password: unHashedPassword,
				}

				store.EXPECT().
					RegisterUser(gomock.Any(), mockdb.EqRegisterUserParams(arg)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RequestLog{}, nil)

				store.EXPECT().
					CreateResponseLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ResponseLog{}, nil)

			},

			// rate limiter is causing trouble
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},

		{
			name: "no email -> bad request",
			body: registerUserRequest{
				Email:    "", // no email
				Password: user.Password,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RequestLog{}, nil)

				store.EXPECT().
					CreateResponseLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ResponseLog{}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				require.Contains(
					t,
					string(data),
					"Field validation for 'Email' failed on the 'required'",
				)
			},
		},

		{
			name: "bad email -> bad request",
			body: registerUserRequest{
				Email:    "thisisabadmail", // bad email
				Password: user.Password,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(user, nil)

				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RequestLog{}, nil)

				store.EXPECT().
					CreateResponseLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ResponseLog{}, nil)

			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "internal server error",
			body: registerUserRequest{
				Email:    user.Email,
				Password: user.Password,
			},

			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrUniqueViolation)

				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RequestLog{}, nil)

				store.EXPECT().
					CreateResponseLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ResponseLog{}, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Marshall body data to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// requireBodyMatchUser checks that the user in the body matches the recieved user
func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	gotUser.ID = user.ID
	require.NoError(t, err)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Password, gotUser.Password)
}

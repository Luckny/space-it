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

func TestRegisterUserAPI(t *testing.T) {
	user, unHashedPassword := mockdb.RandomUser(t)

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
					Email:    user.Email,
					Password: unHashedPassword,
				}

				store.EXPECT().
					RegisterUser(gomock.Any(), mockdb.EqRegisterUserParams(arg)).
					Times(1).
					Return(user, nil)
			},

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
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
	}

	// Run test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init gomock
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// api server with mock store
			server := NewServer(store, config.Config{})
			router := gin.Default()
			router.POST("/users", server.registerUser)

			// request params
			jsonBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			// create request
			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
			require.NoError(t, err)

			// test recorder
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(recorder)
		})
	}

}

func TestLoginUserAPI(t *testing.T) {
	user, _ := mockdb.RandomUser(t)

	testCases := []struct {
		name          string
		setContext    bool
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "should login user",
			setContext: true,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:       "error -> user not in context",
			setContext: false,
		},
	}

	// Run test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init gomock
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			// tc.buildStubs(store)

			// api server with mock store
			server := NewServer(store, config.Config{})
			router := gin.Default()

			if tc.setContext {
				router.Use(func(c *gin.Context) {
					c.Set("user", &user)
				})
			}

			router.POST("/users/login", server.loginUser)

			// create request
			request, err := http.NewRequest(http.MethodPost, "/users/login", nil)
			require.NoError(t, err)

			// test recorder
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			if !tc.setContext {
				require.Panics(t, func() {
					tc.checkResponse(recorder)
				})

			} else {
				// check response
				tc.checkResponse(recorder)
			}
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

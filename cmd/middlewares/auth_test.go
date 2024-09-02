package middlewares

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Luckny/space-it/db/mock"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticate(t *testing.T) {
	user, unHashedPassword := mockdb.RandomUser(t)

	testCases := []struct {
		name          string
		setHeader     bool
		username      string
		password      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "empty params -> Bad request",
			setHeader: true,
			username:  "",
			password:  "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},

		{
			name:      "no header -> no context",
			setHeader: false,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var responseBodyUser db.User
				err = json.Unmarshal(data, &responseBodyUser)
				require.NoError(t, err)

				require.Zero(t, responseBodyUser)
			},
		},

		{
			name:      "Ok",
			setHeader: true,
			username:  user.Email,
			password:  unHashedPassword,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var responseBodyUser db.User
				err = json.Unmarshal(data, &responseBodyUser)
				require.NoError(t, err)

				require.Equal(t, responseBodyUser.Email, user.Email)
				require.Equal(t, responseBodyUser.ID, user.ID)
				require.Equal(t, responseBodyUser.Password, user.Password)
			},
		},

		{
			name:      "not found",
			setHeader: true,
			username:  user.Email,
			password:  unHashedPassword,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},

		{
			name:      "internal error",
			setHeader: true,
			username:  user.Email,
			password:  unHashedPassword,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// server := NewServer(store)
			router := gin.Default()
			router.Use(Authenticate(store))

			router.GET(
				"/getpath",
				func(ctx *gin.Context) {
					user, _ := ctx.Get("user")
					ctx.JSON(http.StatusOK, user)
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/getpath", nil)
			require.NoError(t, err)

			if tc.setHeader {
				request.SetBasicAuth(tc.username, tc.password)
			}

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestRequireAuth(t *testing.T) {
	user, _ := mockdb.RandomUser(t)

	testCases := []struct {
		name          string
		authenticate  bool
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "authenticated request -> ok",
			authenticate: true,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:         "unauthenticated request -> unauthorized",
			authenticate: false,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// api server with mock store
			router := gin.Default()
			if tc.authenticate {
				router.Use(func(c *gin.Context) {
					c.Set("user", &user)
					c.Next()
				})
			}
			router.Use(RequireAuthentication())

			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, nil)
			})

			// create request
			request, err := http.NewRequest(http.MethodGet, "/test", nil)
			require.NoError(t, err)

			// test recorder
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(recorder)
		})
	}
}

package api

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

func TestEnsureJSONContentMiddleware(t *testing.T) {
	testCases := []struct {
		name              string
		method            string
		path              string
		contentTypeHeader string
		checkResponse     func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "Get request -> Ok",
			method: http.MethodGet,
			path:   "/getpath",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},

		{
			name:   "No content type header -> unsupported media type",
			method: http.MethodPost,
			path:   "/postpath",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)
			},
		},

		{
			name:              "with content type header -> ok",
			method:            http.MethodPost,
			contentTypeHeader: "application/json",
			path:              "/postpath",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(nil)
			server.router.Use(EnsureJSONContentType())
			server.router.GET(
				"/getpath",
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			server.router.POST(
				"/postpath",
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(tc.method, tc.path, nil)
			request.Header.Set("Content-Type", tc.contentTypeHeader)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestRateLimiter(t *testing.T) {
	server := NewServer(nil)
	server.router.Use(RateGuard(server.limiter))

	server.router.GET(
		"/getpath",
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	errs := make(chan error)
	responseCode := make(chan int)

	n := int(server.limiter.Limit()) + 1 // number of allowed request + 1

	// n concurrent calls the an enpoint
	for i := 0; i < n; i++ {
		go func() {

			request, err := http.NewRequest(http.MethodGet, "/getpath", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			errs <- err
			responseCode <- recorder.Code
		}()
	}

	codes := []int{}

	// check results
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		code := <-responseCode // put in codes
		codes = append(codes, code)

	}

	// Expect n-1 200 OK responses and 1 429 Too Many Requests response
	expectedCodes := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		expectedCodes[i] = http.StatusOK
	}
	expectedCodes = append(expectedCodes, http.StatusTooManyRequests)
	require.ElementsMatch(t, codes, expectedCodes)

}

func TestAuthenticate(t *testing.T) {
	user, unHashedPassword := randomUser(t)

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

				// Unmarshal the response body into a map
				var responseBody map[string]interface{}
				err = json.Unmarshal(data, &responseBody)
				require.NoError(t, err)

				// Expected response
				expectedBody := map[string]interface{}{
					"user": "",
				}

				require.Equal(t, expectedBody, responseBody)
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

				// Unmarshal the response body into a map
				var responseBody map[string]interface{}
				err = json.Unmarshal(data, &responseBody)
				require.NoError(t, err)

				// Expected response
				expectedBody := map[string]interface{}{
					"user": user.Email,
				}

				require.Equal(t, expectedBody, responseBody)
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

			server := NewServer(store)
			server.router.Use(server.Authenticate())

			server.router.GET(
				"/getpath",
				func(ctx *gin.Context) {
					email, _ := ctx.Get("user")
					ctx.JSON(http.StatusOK, gin.H{"user": email})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/getpath", nil)
			require.NoError(t, err)

			if tc.setHeader {
				request.SetBasicAuth(tc.username, tc.password)
			}

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

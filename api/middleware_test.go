package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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

package middlewares

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
			router := gin.Default()
			router.Use(EnsureJSONContentType())
			router.GET(
				"/getpath",
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			router.POST(
				"/postpath",
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(tc.method, tc.path, nil)
			request.Header.Set("Content-Type", tc.contentTypeHeader)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

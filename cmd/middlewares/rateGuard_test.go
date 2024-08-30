package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	router := gin.Default()
	limiter := rate.NewLimiter(rate.Limit(2), 2)
	router.Use(RateGuard(limiter))

	router.GET(
		"/getpath",
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	errs := make(chan error)
	responseCode := make(chan int)

	n := int(limiter.Limit()) + 1 // number of allowed request + 1

	// n concurrent calls the an enpoint
	for i := 0; i < n; i++ {
		go func() {

			request, err := http.NewRequest(http.MethodGet, "/getpath", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
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

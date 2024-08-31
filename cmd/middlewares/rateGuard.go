package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateGuard(limiter *rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !limiter.Allow() {
			httpx.WriteError(ctx, http.StatusTooManyRequests, fmt.Errorf("too many requests"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

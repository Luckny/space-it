package api

import (
	"fmt"
	"net/http"
	"strings"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// EnsureJSONContentType ensures that post requests content type is application/json
func EnsureJSONContentType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodPost {
			ctx.Next()
			return
		}

		if ctx.GetHeader("Content-Type") != "application/json" {
			writeError(
				ctx,
				http.StatusUnsupportedMediaType,
				fmt.Errorf("Content-Type must be 'application/json'"),
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}

func RateGuard(limiter *rate.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if !limiter.Allow() {
			writeError(ctx, http.StatusTooManyRequests, fmt.Errorf("too many requests"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (server *Server) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		ctx.Set("user", "") // empty user context

		// request is not authenticated
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic") {
			ctx.Next()
			return
		}

		email, password, err := util.ExtractAuthHeader(authHeader)
		if err != nil {
			writeError(ctx, http.StatusBadRequest, err)
			ctx.Abort()
			return
		}

		user, err := server.store.GetUserByEmail(ctx, email)
		if err != nil {
			if err == db.ErrRecordNotFound {
				writeError(ctx, http.StatusNotFound, fmt.Errorf("user not found"))
				return
			}
			writeError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		// set user email in context if no error
		if err := util.CheckPassword(password, user.Password); err == nil {
			ctx.Set("user", user.Email)
		}

		ctx.Next()
	}
}

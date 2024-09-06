package middlewares

import (
	"github.com/Luckny/space-it/pkg/token"
	"github.com/gin-gonic/gin"
)

// Verify token
func VerifyToken(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		tokenID := ctx.GetHeader("X-CSRF-Token")
		if tokenID == "" {
			ctx.Next()
			return
		}

		// validate token
		token, err := maker.VerifyToken(ctx, tokenID)
		if err != nil {
			// Token is invalid
			ctx.Next()
			return
		}

		// Token is valid
		ctx.Set("user", &token.User)
		ctx.Next()

	}
}

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

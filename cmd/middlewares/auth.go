package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/writer"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
)

func Authenticate(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		ctx.Set("user", nil) // empty user context

		// request is not authenticated
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic") {
			ctx.Next()
			return
		}

		email, password, err := util.ExtractAuthHeader(authHeader)
		if err != nil {
			writer.WriteError(ctx, http.StatusBadRequest, err)
			ctx.Abort()
			return
		}

		user, err := store.GetUserByEmail(ctx, email)
		if err != nil {
			if err == db.ErrRecordNotFound {
				writer.WriteError(ctx, http.StatusNotFound, fmt.Errorf("user not found"))
				return
			}
			writer.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		// set user email in context if no error
		if err := util.CheckPassword(password, user.Password); err == nil {
			ctx.Set("user", user)
		}

		ctx.Next()
	}
}

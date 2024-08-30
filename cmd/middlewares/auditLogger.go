package middlewares

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/writer"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuditLogger is middleware that logs incoming HTTP requests and their corresponding responses.
func AuditLogger(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqLogID, err := logRequest(ctx, store)
		if err != nil {
			writer.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		// Wrap the response writer
		w := writer.NewResponseWriter(ctx.Writer)
		ctx.Writer = w

		ctx.Next()

		// log the response
		arg := db.CreateResponseLogParams{ID: reqLogID, Status: int32(w.Status())}
		if err := store.CreateResponseLog(ctx, arg); err != nil {
			writer.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

	}
}

// logRequest logs an incoming HTTP request to the database
func logRequest(ctx *gin.Context, store db.Store) (uuid.UUID, error) {
	u, _ := ctx.Get("user")
	if u != nil {
		// Handle authenticated request logging
		arg := db.CreateAuthenticatedRequestLogParams{
			ID:     util.GenUUID(),
			Path:   ctx.Request.URL.Path,
			Method: ctx.Request.Method,
			UserID: u.(db.User).ID,
		}
		return store.CreateAuthenticatedRequestLog(ctx, arg)
	} else {
		// Handle unauthenticated request logging
		arg := db.CreateUnauthenticatedRequestLogParams{
			ID:     util.GenUUID(),
			Path:   ctx.Request.URL.Path,
			Method: ctx.Request.Method,
		}
		return store.CreateUnauthenticatedRequestLog(ctx, arg)
	}
}

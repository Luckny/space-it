package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			httpx.WriteError(ctx, http.StatusBadRequest, err)
			ctx.Abort()
			return
		}

		user, err := store.GetUserByEmail(ctx, email)
		if err != nil {
			if err == db.ErrRecordNotFound {
				httpx.WriteError(ctx, http.StatusNotFound, fmt.Errorf("user not found"))
				ctx.Abort()
				return
			}
			httpx.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		// set user email in context if no error
		if err := util.CheckPassword(password, user.Password); err == nil {
			ctx.Set("user", &user)
		}

		ctx.Next()
	}
}

func RequireAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := httpx.GetUserFromContext(ctx)
		if err != nil {
			ctx.Header("WWW-Authenticate", "Basic realm=\"/\", charset\"UTF-8\"")
			httpx.WriteError(ctx, http.StatusUnauthorized, fmt.Errorf("who are you?"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func RequireSpacePermission(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, _ := httpx.GetUserFromContext(ctx)

		spaceID, err := uuid.Parse(ctx.Param("spaceID"))
		if err != nil {
			util.InfoLog.Println("i really made it", spaceID)
			httpx.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		arg := db.GetPermissionsByUserAndSpaceIDParams{
			UserID:  user.ID,
			SpaceID: spaceID,
		}

		permission, err := store.GetPermissionsByUserAndSpaceID(ctx, arg)
		if err != nil {
			if err == db.ErrRecordNotFound {
				httpx.WriteError(
					ctx,
					http.StatusForbidden,
					fmt.Errorf("access denied: insufficient permissions"),
				)
				ctx.Abort()
				return
			}
			httpx.WriteError(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		method := ctx.Request.Method
		switch method {
		case "POST":
			if !permission.WritePermission {
				httpx.WriteError(
					ctx,
					http.StatusForbidden,
					fmt.Errorf("insufficient permissions: write access required"),
				)
				ctx.Abort()
				return
			}
			break

		case "GET":
			if !permission.ReadPermission {
				httpx.WriteError(
					ctx,
					http.StatusForbidden,
					fmt.Errorf("insufficient permissions: read access required"),
				)
				ctx.Abort()
				return
			}
			break

		case "DELETE":
			if !permission.DeletePermission {
				httpx.WriteError(
					ctx,
					http.StatusForbidden,
					fmt.Errorf("insufficient permissions: delete access required"),
				)
				ctx.Abort()
				return
			}
			break
		}

		ctx.Next()
	}
}

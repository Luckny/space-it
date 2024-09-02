package middlewares

import (
	"fmt"
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AccessLvl string

const (
	WriteAccess  AccessLvl = "write"
	ViewAccess             = "read"
	DeleteAccess           = "delete"
	AdminAccess            = "admin"
)

func RequireAccessLvl(accessLvl AccessLvl, store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := httpx.GetUserFromContext(ctx)
		if err != nil {
			// user should be authenticated by the auth middlewares
			util.ErrorLog.Panic(err)
			return
		}

		spaceID, err := uuid.Parse(ctx.Param("spaceID"))
		if err != nil {
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
		}

		var userHasAccess bool = false
		switch accessLvl {
		case ViewAccess:
			userHasAccess = permission.ReadPermission
			break
		case WriteAccess:
			userHasAccess = permission.WritePermission
			break
		case DeleteAccess:
			userHasAccess = permission.DeletePermission
			break
		case AdminAccess:
			userHasAccess = permission.ReadPermission && permission.WritePermission &&
				permission.DeletePermission
			break
		}

		if !userHasAccess {
			httpx.WriteError(
				ctx,
				http.StatusForbidden,
				fmt.Errorf("denied: %s access required", accessLvl),
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// input validator
var ValidAccessLvl validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if accessLvls, ok := fieldLevel.Field().Interface().(map[AccessLvl]bool); ok {
		// check if valid
		for key, _ := range accessLvls {
			switch key {
			case WriteAccess, ViewAccess, DeleteAccess:
				continue
			default:
				return false
			}
		}

		return true
	}

	return false
}

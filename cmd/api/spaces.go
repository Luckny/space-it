package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Luckny/space-it/cmd/middlewares"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type createSpaceRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

func (server *Server) createSpace(ctx *gin.Context) {
	var req createSpaceRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := httpx.GetUserFromContext(ctx)
	if err != nil {
		// user should be authenticated by the auth middlewares
		util.ErrorLog.Panic(err)
		return
	}

	arg := db.CreateSpaceTxParams{
		Name:  req.Name,
		Owner: user.ID,
	}

	// creates space and gives it owner permissions
	spaceTxResult, err := server.store.CreateSpaceTx(ctx, arg)
	if err != nil {
		handleCreateSpaceError(ctx, err)
		return
	}

	httpx.WriteResponse(ctx, http.StatusCreated, spaceTxResult.Space)
}

type addMemberToSpaceRequest struct {
	UserID      uuid.UUID                      `json:"user_id"     binding:"required"`
	Permissions map[middlewares.AccessLvl]bool `json:"permissions" binding:"required,accesslvl"`
}

func (server *Server) addMemberToSpace(ctx *gin.Context) {
	var req addMemberToSpaceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	spaceID, err := uuid.Parse(ctx.Param("spaceID"))
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		ctx.Abort()
		return
	}

	arg := db.CreatePermissionParams{
		UserID:           req.UserID,
		SpaceID:          spaceID,
		WritePermission:  req.Permissions[middlewares.WriteAccess],
		ReadPermission:   req.Permissions[middlewares.ViewAccess],
		DeletePermission: req.Permissions[middlewares.DeleteAccess],
	}
	permission, err := server.store.CreatePermission(ctx, arg)

	httpx.WriteResponse(ctx, http.StatusCreated, permission)
}

func handleCreateSpaceError(ctx *gin.Context, err error) {
	var pgErr *pgconn.PgError
	// if not a pg error return generic error
	if !errors.As(err, &pgErr) {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	switch pgErr.Code {
	case db.ErrUniqueViolation.Code:
		httpx.WriteError(ctx, http.StatusConflict, fmt.Errorf("space with name already exists"))
		break

	default:
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		break
	}

	return

}

package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
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

	// error can safely be ignored here
	// this handler requires authentication
	user, _ := httpx.GetUserFromContext(ctx)

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
		util.InfoLog.Println("im insinge")
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		break
	}

	return

}

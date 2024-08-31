package api

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/gin-gonic/gin"
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

	arg := db.CreateSpaceParams{
		Name:  req.Name,
		Owner: user.ID,
	}

	// TODO: check if space name already exists

	space, err := server.store.CreateSpace(ctx, arg)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	httpx.WriteResponse(ctx, http.StatusCreated, space)
}

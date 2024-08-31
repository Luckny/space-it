package api

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/Luckny/space-it/pkg/httpx"
	"github.com/Luckny/space-it/util"
	"github.com/gin-gonic/gin"
)

type registerUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (server *Server) registerUser(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	// hash the password
	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	arg := db.RegisterUserParams{
		ID:       util.GenUUID(),
		Email:    req.Email,
		Password: passwordHash,
	}

	user, err := server.store.RegisterUser(ctx, arg)
	if err != nil {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	httpx.WriteResponse(ctx, http.StatusCreated, user)
}

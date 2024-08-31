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
		Email:    req.Email,
		Password: passwordHash,
	}

	user, err := server.store.RegisterUser(ctx, arg)
	if err != nil {
		handleRegisterUserError(ctx, err)
		return
	}

	httpx.WriteResponse(ctx, http.StatusCreated, user)
}

func handleRegisterUserError(ctx *gin.Context, err error) {
	var pgErr *pgconn.PgError
	// if not a pg error return generic error
	if !errors.As(err, &pgErr) {
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	switch pgErr.Code {
	case db.ErrUniqueViolation.Code:
		httpx.WriteError(ctx, http.StatusConflict, fmt.Errorf("user already exists"))
		break

	default:
		util.InfoLog.Println("im insinge")
		httpx.WriteError(ctx, http.StatusInternalServerError, err)
		break
	}

	return

}

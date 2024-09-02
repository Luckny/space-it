package token

import "github.com/gin-gonic/gin"

type Maker interface {
	// CreateToken create a new token for a specific user and duration
	CreateToken(ctx *gin.Context, payload *Payload) (string, error)
	// VerifyToken checks if the token is valid
	VerifyToken(ctx *gin.Context, tokenID string) (*Payload, error)
}

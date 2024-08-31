package httpx

import (
	"fmt"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*db.User, error) {
	u, ok := c.Get("user")
	if !ok {
		return nil, fmt.Errorf("error getting user from context")
	}

	user, ok := u.(*db.User)
	if !ok {
		return nil, fmt.Errorf("error getting user from context")
	}

	return user, nil
}

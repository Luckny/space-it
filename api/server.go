package api

import (
	"net/http"

	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	server.router = router
	return server
}

// Run's the api server
func (server *Server) Run(addr string) error {
	return server.router.Run(addr)
}

// writeResponse writes a json response with api required headers
func writeResponse(c *gin.Context, status int, value interface{}) {
	c.Header("Content-Type", "application/json")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "0")
	c.Header("Cache-Control", "no-store")

	c.JSON(status, value)
}

// writeError writes an error response
func writeError(c *gin.Context, status int, err error) {
	if status == http.StatusInternalServerError {
		writeResponse(c, status, map[string]string{"error": "internal server error"})
		return
	}
	writeResponse(c, status, map[string]string{"error": err.Error()})
}
